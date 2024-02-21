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
		Stat        chan Stat
		Conn        *grpc.ClientConn
		disconnects int
	}
)

const (
	CONNECT    Stat = 1
	DISCONNECT Stat = 2
)

func New(node *pb.Node) (*grpc.ClientConn, pb.ClusterRpcServiceClient, <-chan Stat, func() error, error) {
	stat := make(chan Stat, 100)
	var conn *grpc.ClientConn
	var err error
	fmt.Printf("%s:%d\r\n", node.Host, node.Port)
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", node.Host, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(&StatHandler{Stat: stat, Conn: conn}))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	client := pb.NewClusterRpcServiceClient(conn)
	return conn, client, stat, conn.Close, nil
}

func (StatHandler) TagRPC(ctx context.Context, tag *stats.RPCTagInfo) context.Context {
	return ctx
}
func (StatHandler) HandleRPC(context.Context, stats.RPCStats) {

}
func (StatHandler) TagConn(ctx context.Context, tag *stats.ConnTagInfo) context.Context {
	return ctx
}
func (statHandler *StatHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnEnd:
		{
			statHandler.Stat <- DISCONNECT
			statHandler.disconnects += 1
		}
	case *stats.ConnBegin:
		{
			if statHandler.disconnects > 0 {
				statHandler.Stat <- CONNECT
			}
		}
	}
}
