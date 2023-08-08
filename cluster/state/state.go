package state

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vedadiyan/coms/cluster/client"
	pb "github.com/vedadiyan/coms/cluster/proto"
)

var (
	id    string
	mut   sync.RWMutex
	nodes map[string]*pb.Node
	conns map[string]pb.ClusterRpcServiceClient
)

func init() {
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
	go func() {
		for {
			_, err := conn.GetId(context.TODO(), &pb.Void{})
			if err != nil {
				mut.Lock()
				defer mut.Unlock()
				delete(nodes, node.Id)
				delete(conns, node.Id)
				fmt.Println("disconnected")
				break
			}
			<-time.After(time.Second)
		}
	}()
	log.Println("joined", node.Port)
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
	mut.RLock()
	defer mut.RUnlock()
	nodeList := pb.NodeList{}
	for _, value := range nodes {
		nodeList.Nodes = append(nodeList.Nodes, value)
	}
	for key, value := range conns {
		if key == gossiperId {
			continue
		}
		if key == id {
			continue
		}
		value.Gossip(context.TODO(), &nodeList)
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
