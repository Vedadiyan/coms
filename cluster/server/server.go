package server

import (
	"context"
	"encoding/json"
	"log"
	"net"

	pb "github.com/vedadiyan/coms/cluster/proto"
	"github.com/vedadiyan/coms/cluster/state"
	"github.com/vedadiyan/coms/socket"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedClusterRpcServiceServer
}

func (server Server) Gossip(ctx context.Context, nodeList *pb.NodeList) (*pb.Void, error) {
	localNodeList := pb.NodeList{
		Id: state.GetId(),
	}
	localNodeList.Nodes = nodeList.Nodes
	nodesAdded := state.AppendNodes(localNodeList.Nodes)
	if nodesAdded > 0 {
		state.GossipAll(nodeList.Id)
		state.Print()
	}
	return &pb.Void{}, nil
}

func (server Server) Exchange(ctx context.Context, exchangeReq *pb.ExchangeReq) (*pb.Void, error) {
	//notImplemeted := func() { panic("not implemented") }
	log.Println("EXCHANGE")
	switch exchangeReq.Event {
	case "emit:room":
		{
			mapper := make(map[string]string)
			err := json.Unmarshal(exchangeReq.Arg, &mapper)
			if err != nil {
				log.Println(err.Error())
				return nil, err
			}
			socket.SendToRoom(nil, mapper["room"], mapper["message"])
		}
	}
	return &pb.Void{}, nil
}

func (server Server) GetId(ctx context.Context, _ *pb.Void) (*pb.Id, error) {
	id := pb.Id{}
	id.Id = state.GetId()
	return &id, nil
}

func New(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClusterRpcServiceServer(s, &Server{})
	log.Printf("cluster listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Solicit(nodeList *pb.NodeList) {
	for _, node := range nodeList.Nodes {
		state.JoinNode(node)
	}
	state.GossipAll(state.GetId())
}
