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
	if node.Id == id {
		return
	}
	mut.Lock()
	defer mut.Unlock()
	if _, ok := nodes[node.Id]; ok {
		return
	}
	conn, err := client.New(node)
	if err != nil {
		log.Println(err)
		return
	}
	nodes[node.Id] = node
	conns[node.Id] = conn
	go handleDisconnect(conn, node.Id)
	log.Println("joined", node.Port, node.Id)
}

func handleDisconnect(conn pb.ClusterRpcServiceClient, id string) {
	disconnected := false
	mut.RLock()
	localConn := conns[id]
	localNode := nodes[id]
	mut.RUnlock()
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		default:
			{
				newId, err := conn.GetId(context.TODO(), &pb.Void{})
				if err != nil {
					disconnected = true
					mut.Lock()
					delete(nodes, id)
					delete(conns, id)
					mut.Unlock()
					// log.Println("connection lost", id)
				}
				if disconnected && newId != nil {
					mut.Lock()
					localNode.Id = newId.Id
					nodes[newId.Id] = localNode
					conns[newId.Id] = localConn
					mut.Unlock()
					nodeList := pb.NodeList{}
					nodeList.Id = GetId()
					mut.RLock()
					for key, value := range nodes {
						if key == newId.Id {
							continue
						}
						nodeList.Nodes = append(nodeList.Nodes, value)
					}
					mut.RUnlock()
					_, err := conns[newId.Id].Gossip(context.TODO(), &nodeList)
					if err == nil {
						disconnected = false
					}
				}
				<-time.After(time.Millisecond * 25)
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

// func GetNodes(cb func(node *pb.Node, conn pb.ClusterRpcServiceClient)) {
// 	// mut.RLock()
// 	// defer mut.RUnlock()
// 	for key, value := range nodes {
// 		cb(value, conns[key])
// 	}
// }

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
		conn := value
		if key == gossiperId {
			continue
		}
		if key == id {
			continue
		}
		cb = append(cb, func() {
			conn.Gossip(context.TODO(), &nodeList)
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
