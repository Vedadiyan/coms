// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: cluster/proto/cluster-rpc.proto

package cluster

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ClusterRpcServiceClient is the client API for ClusterRpcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterRpcServiceClient interface {
	Gossip(ctx context.Context, in *NodeList, opts ...grpc.CallOption) (*Void, error)
	GetId(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Id, error)
	Exchange(ctx context.Context, in *ExchangeReq, opts ...grpc.CallOption) (*Void, error)
}

type clusterRpcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterRpcServiceClient(cc grpc.ClientConnInterface) ClusterRpcServiceClient {
	return &clusterRpcServiceClient{cc}
}

func (c *clusterRpcServiceClient) Gossip(ctx context.Context, in *NodeList, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/coms.proto.cluster.ClusterRpcService/Gossip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterRpcServiceClient) GetId(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Id, error) {
	out := new(Id)
	err := c.cc.Invoke(ctx, "/coms.proto.cluster.ClusterRpcService/GetId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterRpcServiceClient) Exchange(ctx context.Context, in *ExchangeReq, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/coms.proto.cluster.ClusterRpcService/Exchange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterRpcServiceServer is the server API for ClusterRpcService service.
// All implementations must embed UnimplementedClusterRpcServiceServer
// for forward compatibility
type ClusterRpcServiceServer interface {
	Gossip(context.Context, *NodeList) (*Void, error)
	GetId(context.Context, *Void) (*Id, error)
	Exchange(context.Context, *ExchangeReq) (*Void, error)
	mustEmbedUnimplementedClusterRpcServiceServer()
}

// UnimplementedClusterRpcServiceServer must be embedded to have forward compatible implementations.
type UnimplementedClusterRpcServiceServer struct {
}

func (UnimplementedClusterRpcServiceServer) Gossip(context.Context, *NodeList) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Gossip not implemented")
}
func (UnimplementedClusterRpcServiceServer) GetId(context.Context, *Void) (*Id, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetId not implemented")
}
func (UnimplementedClusterRpcServiceServer) Exchange(context.Context, *ExchangeReq) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exchange not implemented")
}
func (UnimplementedClusterRpcServiceServer) mustEmbedUnimplementedClusterRpcServiceServer() {}

// UnsafeClusterRpcServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterRpcServiceServer will
// result in compilation errors.
type UnsafeClusterRpcServiceServer interface {
	mustEmbedUnimplementedClusterRpcServiceServer()
}

func RegisterClusterRpcServiceServer(s grpc.ServiceRegistrar, srv ClusterRpcServiceServer) {
	s.RegisterService(&ClusterRpcService_ServiceDesc, srv)
}

func _ClusterRpcService_Gossip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterRpcServiceServer).Gossip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coms.proto.cluster.ClusterRpcService/Gossip",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterRpcServiceServer).Gossip(ctx, req.(*NodeList))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterRpcService_GetId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterRpcServiceServer).GetId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coms.proto.cluster.ClusterRpcService/GetId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterRpcServiceServer).GetId(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterRpcService_Exchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExchangeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterRpcServiceServer).Exchange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coms.proto.cluster.ClusterRpcService/Exchange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterRpcServiceServer).Exchange(ctx, req.(*ExchangeReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ClusterRpcService_ServiceDesc is the grpc.ServiceDesc for ClusterRpcService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterRpcService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "coms.proto.cluster.ClusterRpcService",
	HandlerType: (*ClusterRpcServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Gossip",
			Handler:    _ClusterRpcService_Gossip_Handler,
		},
		{
			MethodName: "GetId",
			Handler:    _ClusterRpcService_GetId_Handler,
		},
		{
			MethodName: "Exchange",
			Handler:    _ClusterRpcService_Exchange_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster/proto/cluster-rpc.proto",
}
