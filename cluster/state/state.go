package state

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vedadiyan/coms/cluster/client"
	pb "github.com/vedadiyan/coms/cluster/proto"
)

var (
	ctx   context.Context
	id    string
	mut   sync.RWMutex
	nodes map[string]*pb.Node
	conns map[string]pb.ClusterRpcServiceClient
)

func init() {
	ctx = context.TODO()
	nodes = make(map[string]*pb.Node)
	conns = make(map[string]pb.ClusterRpcServiceClient)
	id = uuid.New().String()
}

func JoinNode(node *pb.Node) {
	conn, err := client.New(node)
	if err != nil {
		log.Println(err)
		return
	}
	mut.Lock()
	defer mut.Unlock()
	nodes[node.Id] = node
	conns[node.Id] = conn
	go handleDisconnect(conn, node.Id)
	log.Println("joined", node.Port)
}

func handleDisconnect(conn pb.ClusterRpcServiceClient, id string) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		default:
			{
				_, err := conn.GetId(context.TODO(), &pb.Void{})
				if err != nil {
					mut.Lock()
					delete(nodes, id)
					delete(conns, id)
					mut.Unlock()
					log.Println("connection lost", id)
					return
				}
				<-time.After(time.Second)
			}
		}
	}
}

func JoinSelf(node *pb.Node) {
	mut.Lock()
	defer mut.Unlock()
	nodes[node.Id] = node
	conns[node.Id] = nil
}

func GetNodes(cb func(node *pb.Node, conn pb.ClusterRpcServiceClient)) {
	// mut.RLock()
	// defer mut.RUnlock()
	for key, value := range nodes {
		cb(value, conns[key])
	}
}

func AppendNodes(_nodes []*pb.Node) int {
	mut.RLock()
	local := make([]*pb.Node, 0)
	for _, node := range _nodes {
		if _, ok := nodes[node.Id]; ok {
			continue
		}
		local = append(local, node)
	}
	mut.RUnlock()
	for _, node := range local {
		JoinNode(node)
	}
	return len(local)
}

func GossipAll(gossiperId string) {
	cb := make([]func(), 0)
	mut.RLock()
	nodeList := pb.NodeList{}
	for _, value := range nodes {
		nodeList.Nodes = append(nodeList.Nodes, value)
	}
	for key, value := range conns {
		vv := value
		if key == gossiperId {
			continue
		}
		if key == id {
			continue
		}
		cb = append(cb, func() {
			vv.Gossip(context.TODO(), &nodeList)
		})
	}
	mut.RUnlock()
	for _, fn := range cb {
		fn()
	}
}

func GetId() string {
	return id
}

func Print() {
	mut.RLock()
	defer mut.RUnlock()
	log.Println("nodes connected", len(nodes)-1)
}
