package socket

import (
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

type Options struct {
	authenticate func(r *http.Request) (bool, error)
	intercept    func(socket *Socket, message []byte) (bool, error)
	closeHandler func(socket *Socket)
}

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
	_mut     sync.RWMutex
	_sockets map[string]*Socket
	_rooms   map[string]map[string]*Socket
	_options Options
)

func init() {
	_sockets = make(map[string]*Socket)
	_rooms = make(map[string]map[string]*Socket)
}

func New(host string, hub string, options ...func(option *Options)) {
	for _, option := range options {
		option(&_options)
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc(hub, func(w http.ResponseWriter, r *http.Request) {
		if _options.authenticate != nil {
			next, err := _options.authenticate(r)
			if err != nil {
				log.Println(err.Error())
				return
			}
			if !next {
				return
			}
		}
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
		_mut.Lock()
		_sockets[id] = socket
		_mut.Unlock()
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
		if _options.intercept != nil {
			next, err := _options.intercept(socket, message)
			if err != nil {
				log.Println(err.Error())
				return
			}
			if !next {
				return
			}
		}
		go func() {
			log.Println(string(message))
			data := pb.ExchangeReq{}
			err = proto.Unmarshal(message, &data)
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
	localRooms := make([]string, 0)
	_mut.RLock()
	for key := range _rooms {
		localRooms = append(localRooms, key)
	}
	_mut.RUnlock()
	_mut.Lock()
	delete(_sockets, socket.id)
	for _, room := range localRooms {
		delete(_rooms[room], socket.id)
	}
	_mut.Unlock()
	if _options.closeHandler != nil {
		_options.closeHandler(socket)
	}
}

func (socket *Socket) JoinRoom(room string) {
	_mut.Lock()
	defer _mut.Unlock()
	if _, ok := _rooms[room]; !ok {
		_rooms[room] = make(map[string]*Socket)
	}
	_rooms[room][socket.id] = socket
	socket.sendToRoom(common.ROOM_JOIN, room, []byte("Hello!"))
}

func (socket *Socket) LeaveRoom(room string) {
	_mut.Lock()
	defer _mut.Unlock()
	if _, ok := _rooms[room]; !ok {
		return
	}
	delete(_rooms[room], socket.id)
	if len(_rooms[room]) == 0 {
		delete(_rooms, room)
	}
	socket.sendToRoom(common.ROOM_LEAVE, room, []byte("Goodbye!"))
}

func (socket *Socket) SendToRoom(room string, message []byte) {
	socket.sendToRoom(common.EMIT_ROOM, room, message)
}

func (socket *Socket) sendToRoom(event common.Events, room string, message []byte) {
	msg := pb.ExchangeReq{
		Event:   string(event),
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
	_mut.RLock()
	defer _mut.RUnlock()
	for _, sock := range _rooms[room] {
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
	_mut.RLock()
	defer _mut.RUnlock()
	sock, ok := _sockets[to]
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
	_mut.RLock()
	defer _mut.RUnlock()
	for _, sock := range _rooms[msg.To] {
		go sock.Emit(bytes)
	}
}

func Send(msg *pb.ExchangeReq) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_mut.RLock()
	defer _mut.RUnlock()
	sock, ok := _sockets[msg.To]
	if !ok {
		return
	}
	go sock.Emit(bytes)
}

func WithAuthentication(authenticator func(r *http.Request) (bool, error)) func(option *Options) {
	return func(option *Options) {
		option.authenticate = authenticator
	}
}

func WithInterceptor(interceptor func(socket *Socket, message []byte) (bool, error)) func(option *Options) {
	return func(option *Options) {
		option.intercept = interceptor
	}
}

func WithCloseHandler(closeHandler func(socket *Socket)) func(option *Options) {
	return func(option *Options) {
		option.closeHandler = closeHandler
	}
}
