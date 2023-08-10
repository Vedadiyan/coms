package socket

import (
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
	conn *websocket.Conn
}

func (socket *Socket) Emit(event string, payload map[string]any) {
	data := map[string]any{
		"event":     event,
		"timestamp": time.Now(),
		"data":      payload,
	}
	socket.conn.WriteJSON(data)
}

var (
	mut     sync.RWMutex
	sockets map[string]*Socket
)

func init() {
	sockets = make(map[string]*Socket)
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
		socket.conn = conn
		mut.Lock()
		sockets[id] = socket
		mut.Unlock()
		go socketHandler(socket)
	})
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
		fmt.Println(string(message))
	}
}
