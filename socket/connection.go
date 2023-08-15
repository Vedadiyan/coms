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

type Options struct {
	authenticate  func(r *http.Request) (bool, error)
	intercept     func(socket *Socket, message []byte) (bool, error)
	closeHandler  func(socket *Socket)
	invisibleMode func(r *http.Request) (bool, error)
}

type Socket struct {
	id            string
	conn          *websocket.Conn
	header        http.Header
	mut           sync.Mutex
	invisibleMode bool
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
	_groups  map[string]map[string]*Socket
	_options Options
)

func init() {
	_sockets = make(map[string]*Socket)
	_groups = make(map[string]map[string]*Socket)
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
		socket := &Socket{}
		if _options.invisibleMode != nil {
			invisibleMode, err := _options.invisibleMode(r)
			if err != nil {
				log.Println(err.Error())
				return
			}
			socket.invisibleMode = invisibleMode
		}
		id := uuid.New().String()
		fmt.Println(id)
		header := http.Header{}
		header.Add("id", id)
		conn, err := upgrader.Upgrade(w, r, header)
		if err != nil {
			panic(err)
		}
		socket.id = id
		socket.conn = conn
		socket.header = r.Header
		_mut.Lock()
		_sockets[id] = socket
		_mut.Unlock()
		go socketHandler(socket)
	})
	http.HandleFunc("/monit", func(w http.ResponseWriter, r *http.Request) {
		_mut.RLock()
		sockets := make([]string, 0)
		for key := range _sockets {
			sockets = append(sockets, key)
		}
		groups := make(map[string][]string)
		for key := range _groups {
			sockets := make([]string, 0)
			for key := range _groups[key] {
				sockets = append(sockets, key)
			}
			groups[key] = sockets
		}
		_mut.RUnlock()
		output := map[string]any{
			"sockets": sockets,
			"groups":  groups,
		}
		json, err := json.Marshal(output)
		if err != nil {
			w.Header().Add("tatus", "500")
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("content-type", "application/json")
		w.Write(json)
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
			case common.GROUP_JOIN:
				{
					socket.JoinGroup(data.To)
				}
			case common.GROUP_LEAVE:
				{
					socket.LeaveGroup(data.To)
				}
			case common.EMIT_GROUP:
				{
					socket.SendToGroup(data.To, data.Message)
				}
			case common.EMIT_SOCKET:
				{
					socket.Send(data.To, data.Message)
				}
			}
		}()
	}
	localGroups := make([]string, 0)
	_mut.RLock()
	for key := range _groups {
		localGroups = append(localGroups, key)
	}
	_mut.RUnlock()
	_mut.Lock()
	delete(_sockets, socket.id)
	for _, group := range localGroups {
		delete(_groups[group], socket.id)
	}
	_mut.Unlock()
	if _options.closeHandler != nil {
		_options.closeHandler(socket)
	}
}

func (socket *Socket) JoinGroup(group string) {
	_mut.Lock()
	if _, ok := _groups[group]; !ok {
		_groups[group] = make(map[string]*Socket)
	}
	_groups[group][socket.id] = socket
	_mut.Unlock()
	if !socket.invisibleMode {
		socket.sendToGroup(common.GROUP_JOIN, group, []byte("Hello!"))
	}
}

func (socket *Socket) LeaveGroup(group string) {
	_mut.RLock()
	if _, ok := _groups[group]; !ok {
		_mut.RUnlock()
		return
	}
	_mut.RUnlock()
	_mut.Lock()
	delete(_groups[group], socket.id)
	if len(_groups[group]) == 0 {
		delete(_groups, group)
	}
	_mut.Unlock()
	if !socket.invisibleMode {
		socket.sendToGroup(common.GROUP_LEAVE, group, []byte("Goodbye!"))
	}
}

func (socket *Socket) SendToGroup(group string, message []byte) {
	socket.sendToGroup(common.EMIT_GROUP, group, message)
}

func (socket *Socket) sendToGroup(event common.Events, group string, message []byte) {
	msg := pb.ExchangeReq{
		Event:   string(event),
		From:    state.GetId(),
		To:      group,
		Message: message,
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	go state.ExchangeAll(&msg)
	_mut.RLock()
	for _, sock := range _groups[group] {
		if sock == socket {
			continue
		}
		go sock.Emit(bytes)
	}
	_mut.RUnlock()
}

func (socket *Socket) Send(to string, message []byte) {
	msg := pb.ExchangeReq{
		Event:   string(common.EMIT_GROUP),
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

func SendToGroup(msg *pb.ExchangeReq) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_mut.RLock()
	defer _mut.RUnlock()
	for _, sock := range _groups[msg.To] {
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

func WithInvisibleMode(invisibleMode func(r *http.Request) (bool, error)) func(option *Options) {
	return func(option *Options) {
		option.invisibleMode = invisibleMode
	}
}
