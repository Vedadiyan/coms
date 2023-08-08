package client

import (
	"fmt"

	pb "github.com/vedadiyan/coms/cluster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(node *pb.Node) (pb.ClusterRpcServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.Host, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewClusterRpcServiceClient(conn), nil
}
