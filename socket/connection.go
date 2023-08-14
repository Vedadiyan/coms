package socket

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vedadiyan/coms/cluster/state"
)

type Socket struct {
	id     string
	conn   *websocket.Conn
	header http.Header
	mut    sync.Mutex
}

func (socket *Socket) Emit(event string, payload map[string]any) {
	data := map[string]any{
		"event":     event,
		"timestamp": time.Now(),
		"data":      payload,
	}
	socket.mut.Lock()
	defer socket.mut.Unlock()
	socket.conn.WriteJSON(data)
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
		state.ExchangeAll("socket:connected", []byte(id))
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
			data := make(map[string]string)
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Println(err.Error())
			}
			switch data["event"] {
			case "room:join":
				{
					socket.JoinRoom(data["room"])
				}
			case "room:leave":
				{
					socket.LeaveRoom(data["room"])
				}
			case "emit:room":
				{
					socket.SendToRoom(data["room"], data["message"])
				}
			case "emnit:socket":
				{
					socket.Send(data["to"], data["message"])
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

func (socket *Socket) SendToRoom(room string, message string) {
	msg := map[string]any{
		"from":    "",
		"room":    room,
		"message": message,
	}
	json, err := json.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if socket != nil {
		go state.ExchangeAll("emit:room", json)
	}
	mut.RLock()
	defer mut.RUnlock()
	for _, sock := range rooms[room] {
		if sock == socket {
			continue
		}
		go sock.Emit("message", msg)
	}
}

func (socket *Socket) Send(to string, message string) {
	msg := map[string]any{
		"from":    "",
		"message": message,
	}
	json, err := json.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if socket != nil {
		go state.ExchangeAll("emit:socket", json)
	}
	mut.RLock()
	defer mut.RUnlock()
	sock, ok := sockets[to]
	if !ok {
		return
	}
	go sock.Emit("message", msg)
}

func SendToRoom(room string, message map[string]any) {
	mut.RLock()
	defer mut.RUnlock()
	for _, sock := range rooms[room] {
		go sock.Emit("message", message)
	}
}

func Send(to string, message map[string]any) {
	mut.RLock()
	defer mut.RUnlock()
	sock, ok := sockets[to]
	if !ok {
		return
	}
	go sock.Emit("message", message)
}
