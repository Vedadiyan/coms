package server

import (
	"context"
	"log"
	"net"

	pb "github.com/vedadiyan/coms/cluster/proto"
	"github.com/vedadiyan/coms/cluster/state"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedClusterRpcServiceServer
}

func (server Server) Gossip(ctx context.Context, nodeList *pb.NodeList) (*pb.Void, error) {
	localNodeList := pb.NodeList{
		Id: state.GetId(),
	}
	for _, node := range nodeList.Nodes {
		if node.Id == state.GetId() {
			continue
		}
		localNodeList.Nodes = append(localNodeList.Nodes, node)
	}
	nodesAdded := state.AppendNodes(localNodeList.Nodes)
	if nodesAdded > 0 {
		state.GossipAll(nodeList.Id)
		state.Print()
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
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Route(nodeList *pb.NodeList) {
	for _, node := range nodeList.Nodes {
		if node.Id == state.GetId() {
			continue
		}
		state.JoinNode(node)
	}
	state.GetNodes(func(node *pb.Node, conn pb.ClusterRpcServiceClient) {
		if node.Id == state.GetId() {
			return
		}
		conn.Gossip(context.TODO(), nodeList)
	})
}
