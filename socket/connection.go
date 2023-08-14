package socket

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	pb "github.com/vedadiyan/coms/cluster/proto"
	"github.com/vedadiyan/coms/cluster/state"
	"github.com/vedadiyan/coms/common"
	"google.golang.org/protobuf/proto"
)

type Socket struct {
	id     string
	conn   *websocket.Conn
	header http.Header
	mut    sync.Mutex
}

func (socket *Socket) Emit(data []byte) {
	socket.mut.Lock()
	defer socket.mut.Unlock()
	err := socket.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println(err.Error())
	}
}

var (
	mut     sync.RWMutex
	sockets map[string]*Socket
	rooms   map[string]map[string]*Socket
)

func init() {
	sockets = make(map[string]*Socket)
	rooms = make(map[string]map[string]*Socket)
}

func New(host string, hub string) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc(hub, func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		fmt.Println(id)
		header := http.Header{}
		header.Add("id", id)
		conn, err := upgrader.Upgrade(w, r, header)
		if err != nil {
			panic(err)
		}
		socket := &Socket{}
		socket.id = id
		socket.conn = conn
		socket.header = r.Header
		mut.Lock()
		sockets[id] = socket
		mut.Unlock()
		go socketHandler(socket)
	})
	log.Printf("Websocket listening at %s%s\r\n", host, hub)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		panic(err)
	}
}

func socketHandler(socket *Socket) {
	for {
		_, reader, err := socket.conn.NextReader()
		if err != nil {
			break
		}
		message, err := io.ReadAll(reader)
		if err != nil {
			log.Println(err.Error())
		}
		go func() {
			log.Println(string(message))
			data := pb.ExchangeReq{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Println(err.Error())
			}
			switch common.Events(data.Event) {
			case common.ROOM_JOIN:
				{
					socket.JoinRoom(data.To)
				}
			case common.ROOM_LEAVE:
				{
					socket.LeaveRoom(data.To)
				}
			case common.EMIT_ROOM:
				{
					socket.SendToRoom(data.To, data.Message)
				}
			case common.EMIT_SOCKET:
				{
					socket.Send(data.To, data.Message)
				}
			}
		}()
	}
}

func (socket *Socket) JoinRoom(room string) {
	mut.Lock()
	defer mut.Unlock()
	if _, ok := rooms[room]; !ok {
		rooms[room] = make(map[string]*Socket)
	}
	rooms[room][socket.id] = socket
	fmt.Println(socket.id)
}

func (socket *Socket) LeaveRoom(room string) {
	mut.Lock()
	defer mut.Unlock()
	if _, ok := rooms[room]; !ok {
		return
	}
	delete(rooms[room], socket.id)
	if len(rooms[room]) == 0 {
		delete(rooms, room)
	}
}

func (socket *Socket) SendToRoom(room string, message []byte) {
	msg := pb.ExchangeReq{
		Event:   string(common.EMIT_ROOM),
		From:    state.GetId(),
		To:      room,
		Message: message,
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	go state.ExchangeAll(&msg)
	mut.RLock()
	defer mut.RUnlock()
	for _, sock := range rooms[room] {
		if sock == socket {
			continue
		}
		go sock.Emit(bytes)
	}
}

func (socket *Socket) Send(to string, message []byte) {
	msg := pb.ExchangeReq{
		Event:   string(common.EMIT_ROOM),
		From:    state.GetId(),
		To:      to,
		Message: message,
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	go state.ExchangeAll(&msg)
	mut.RLock()
	defer mut.RUnlock()
	sock, ok := sockets[to]
	if !ok {
		return
	}
	go sock.Emit(bytes)
}

func SendToRoom(msg *pb.ExchangeReq) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	mut.RLock()
	defer mut.RUnlock()
	for _, sock := range rooms[msg.To] {
		go sock.Emit(bytes)
	}
}

func Send(msg *pb.ExchangeReq) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	mut.RLock()
	defer mut.RUnlock()
	sock, ok := sockets[msg.To]
	if !ok {
		return
	}
	go sock.Emit(bytes)
}
