package socket

import (
	"context"
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
	intercept     func(socket *Socket, message *pb.ExchangeReq) (bool, error)
	closeHandler  func(socket *Socket)
	invisibleMode func(r *http.Request) (bool, error)
}

type Socket struct {
	id            string
	ip            string
	conn          *websocket.Conn
	header        http.Header
	mut           sync.Mutex
	invisibleMode bool
	claims        map[string]any
}

func (socket *Socket) Emit(data []byte) {
	socket.mut.Lock()
	defer socket.mut.Unlock()
	err := socket.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println(err.Error())
	}
}

func (socket *Socket) Reply(inbox string, status string) {
	msg := pb.ExchangeReq{
		Event:   inbox,
		From:    state.GetId(),
		Message: []byte(status),
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	socket.Emit(bytes)
}

func (socket *Socket) Id() string {
	return socket.id
}

func (socket *Socket) IP() string {
	return socket.ip
}

func (socket *Socket) SetClaim(key string, value any) {
	socket.claims[key] = value
}

func (socket *Socket) GetClaim(key string) (any, bool) {
	value, ok := socket.claims[key]
	return value, ok
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

func New(ctx context.Context, host string, hub string, options ...func(option *Options)) {
	for _, option := range options {
		option(&_options)
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(hub, func(w http.ResponseWriter, r *http.Request) {
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
		socket.ip = r.RemoteAddr
		socket.conn = conn
		socket.header = r.Header
		socket.claims = make(map[string]any)
		_mut.Lock()
		_sockets[id] = socket
		_mut.Unlock()
		go socketHandler(socket)
	})

	serveMux.HandleFunc("/monit", func(w http.ResponseWriter, r *http.Request) {
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
			w.Header().Add("status", "500")
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("content-type", "application/json")
		w.Write(json)
	})
	log.Printf("Websocket listening at %s%s\r\n", host, hub)
	server := http.Server{}
	server.Addr = host
	server.Handler = serveMux
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.TODO())
		if err != nil {
			panic(err)
		}
	}()
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
			data := pb.ExchangeReq{}
			err = proto.Unmarshal(message, &data)
			if err != nil {
				log.Println(err.Error())
			}
			if _options.intercept != nil {
				next, err := _options.intercept(socket, &data)
				if err != nil {
					log.Println(err.Error())
					return
				}
				if !next {
					if data.Reply != nil {
						socket.Reply(*data.Reply, "400")
					}
					return
				}
			}
			switch common.Events(data.Event) {
			case common.JOIN_GROUP:
				{
					socket.JoinGroup(data.To, data.Reply)
				}
			case common.LEAVE_GROUP:
				{
					socket.LeaveGroup(data.To, data.Reply)
				}
			case common.EMIT_GROUP:
				{
					socket.SendToGroup(data.To, data.Reply, data.Message)
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

func (socket *Socket) JoinGroup(group string, inbox *string) {
	_mut.Lock()
	if _, ok := _groups[group]; !ok {
		_groups[group] = make(map[string]*Socket)
	}
	_groups[group][socket.id] = socket
	_mut.Unlock()
	if inbox != nil {
		socket.Reply(*inbox, "200")
	}
	if !socket.invisibleMode {
		socket.sendToGroup(common.JOIN_GROUP, group, nil, []byte("Hello!"))
	}
}

func (socket *Socket) LeaveGroup(group string, inbox *string) {
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
	if inbox != nil {
		socket.Reply(*inbox, "200")
	}
	if !socket.invisibleMode {
		socket.sendToGroup(common.LEAVE_GROUP, group, nil, []byte("Goodbye!"))
	}
}

func (socket *Socket) SendToGroup(group string, inbox *string, message []byte) {
	socket.sendToGroup(common.EMIT_GROUP, group, inbox, message)
}

func (socket *Socket) sendToGroup(event common.Events, group string, inbox *string, message []byte) {
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
	if inbox != nil {
		socket.Reply(*inbox, "200")
	}
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

func WithInterceptor(interceptor func(socket *Socket, message *pb.ExchangeReq) (bool, error)) func(option *Options) {
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
