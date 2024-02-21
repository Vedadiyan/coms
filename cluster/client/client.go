package client

import (
	"context"
	"fmt"

	pb "github.com/vedadiyan/coms/cluster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"
)

type (
	Stat        int
	StatHandler struct {
		Stat chan Stat
	}
)

const (
	CONNECT    Stat = 1
	DISCONNECT Stat = 2
)

func New(node *pb.Node) (pb.ClusterRpcServiceClient, <-chan Stat, func() error, error) {
	stat := make(chan Stat, 1)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.Host, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(StatHandler{Stat: stat}))
	if err != nil {
		return nil, nil, nil, err
	}
	return pb.NewClusterRpcServiceClient(conn), stat, conn.Close, nil
}

func (StatHandler) TagRPC(ctx context.Context, tag *stats.RPCTagInfo) context.Context {
	return ctx
}
func (StatHandler) HandleRPC(context.Context, stats.RPCStats) {

}
func (StatHandler) TagConn(ctx context.Context, tag *stats.ConnTagInfo) context.Context {
	return ctx
}
func (statHandler StatHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnEnd:
		{
			statHandler.Stat <- DISCONNECT
		}
	case *stats.ConnBegin:
		{
			statHandler.Stat <- CONNECT
		}
	}
}
