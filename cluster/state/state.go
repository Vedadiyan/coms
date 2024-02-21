package state

import (
	"context"
	"fmt"
	"log"
	"sync"

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
	if node.Id == id {
		return
	}
	mut.Lock()
	defer mut.Unlock()
	if _, ok := nodes[node.Id]; ok {
		return
	}
	conn, stat, closer, err := client.New(node)
	if err != nil {
		log.Println(err)
		return
	}
	nodes[node.Id] = node
	conns[node.Id] = conn
	go handleDisconnect(conn, stat, closer, node.Id)
	log.Println("joined", node.Port, node.Id)
}

func handleDisconnect(conn pb.ClusterRpcServiceClient, stat <-chan client.Stat, closer func() error, id string) {
	disconnected := false
	mut.RLock()
	localConn := conns[id]
	localNode := nodes[id]
	mut.RUnlock()
	for stat := range stat {
		switch stat {
		case client.DISCONNECT:
			{
				mut.Lock()
				delete(nodes, id)
				delete(conns, id)
				mut.Unlock()
				disconnected = true
				Print()
			}
		case client.CONNECT:
			{
				if !disconnected {
					fmt.Println("skip")
					continue
				}
				newId, err := conn.GetId(context.TODO(), &pb.Void{})
				if err != nil {
					log.Println(err)
					continue
				}
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
				_, err = conns[newId.Id].Gossip(context.TODO(), &nodeList)
				if err != nil {
					log.Println(err)
					continue
				}
				disconnected = false
				Print()
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

func AppendNodes(incomingNodes []*pb.Node) int {
	mut.RLock()
	local := make([]*pb.Node, 0)
	for _, node := range incomingNodes {
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

func ExchangeAll(msg *pb.ExchangeReq) {
	cb := make([]func(), 0)
	mut.RLock()
	nodeList := pb.NodeList{}
	for _, value := range nodes {
		nodeList.Nodes = append(nodeList.Nodes, value)
	}
	for key, value := range conns {
		conn := value
		if key == id {
			continue
		}
		cb = append(cb, func() {
			conn.Exchange(context.TODO(), msg)
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
