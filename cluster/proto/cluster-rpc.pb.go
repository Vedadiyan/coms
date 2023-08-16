// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1-devel
// 	protoc        v4.22.2
// source: cluster/proto/cluster-rpc.proto

package cluster

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Host string `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	Port int32  `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_cluster_proto_cluster_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Node) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Node) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type NodeList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Nodes []*Node `protobuf:"bytes,2,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

func (x *NodeList) Reset() {
	*x = NodeList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeList) ProtoMessage() {}

func (x *NodeList) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeList.ProtoReflect.Descriptor instead.
func (*NodeList) Descriptor() ([]byte, []int) {
	return file_cluster_proto_cluster_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *NodeList) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *NodeList) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

type Void struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Void) Reset() {
	*x = Void{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Void) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Void) ProtoMessage() {}

func (x *Void) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Void.ProtoReflect.Descriptor instead.
func (*Void) Descriptor() ([]byte, []int) {
	return file_cluster_proto_cluster_rpc_proto_rawDescGZIP(), []int{2}
}

type Id struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Id) Reset() {
	*x = Id{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Id) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Id) ProtoMessage() {}

func (x *Id) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Id.ProtoReflect.Descriptor instead.
func (*Id) Descriptor() ([]byte, []int) {
	return file_cluster_proto_cluster_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *Id) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ExchangeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From    string  `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	Event   string  `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
	To      string  `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Reply   *string `protobuf:"bytes,4,opt,name=reply,proto3,oneof" json:"reply,omitempty"`
	Message []byte  `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ExchangeReq) Reset() {
	*x = ExchangeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExchangeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExchangeReq) ProtoMessage() {}

func (x *ExchangeReq) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_proto_cluster_rpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExchangeReq.ProtoReflect.Descriptor instead.
func (*ExchangeReq) Descriptor() ([]byte, []int) {
	return file_cluster_proto_cluster_rpc_proto_rawDescGZIP(), []int{4}
}

func (x *ExchangeReq) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *ExchangeReq) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *ExchangeReq) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *ExchangeReq) GetReply() string {
	if x != nil && x.Reply != nil {
		return *x.Reply
	}
	return ""
}

func (x *ExchangeReq) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_cluster_proto_cluster_rpc_proto protoreflect.FileDescriptor

var file_cluster_proto_cluster_rpc_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x63, 0x6f, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x22, 0x3e, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x4a, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x2e, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x22, 0x06, 0x0a, 0x04, 0x56, 0x6f, 0x69, 0x64, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x86, 0x01, 0x0a, 0x0b, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x19, 0x0a, 0x05, 0x72, 0x65, 0x70,
	0x6c, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x72, 0x65, 0x70, 0x6c,
	0x79, 0x88, 0x01, 0x01, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x08,
	0x0a, 0x06, 0x5f, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x32, 0xd7, 0x01, 0x0a, 0x11, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40,
	0x0a, 0x06, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x12, 0x1c, 0x2e, 0x63, 0x6f, 0x6d, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x6f, 0x69, 0x64,
	0x12, 0x39, 0x0a, 0x05, 0x47, 0x65, 0x74, 0x49, 0x64, 0x12, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x56,
	0x6f, 0x69, 0x64, 0x1a, 0x16, 0x2e, 0x63, 0x6f, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x49, 0x64, 0x12, 0x45, 0x0a, 0x08, 0x45,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1f, 0x2e, 0x63, 0x6f, 0x6d, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x6f,
	0x69, 0x64, 0x42, 0x11, 0x5a, 0x0f, 0x61, 0x75, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cluster_proto_cluster_rpc_proto_rawDescOnce sync.Once
	file_cluster_proto_cluster_rpc_proto_rawDescData = file_cluster_proto_cluster_rpc_proto_rawDesc
)

func file_cluster_proto_cluster_rpc_proto_rawDescGZIP() []byte {
	file_cluster_proto_cluster_rpc_proto_rawDescOnce.Do(func() {
		file_cluster_proto_cluster_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_cluster_proto_cluster_rpc_proto_rawDescData)
	})
	return file_cluster_proto_cluster_rpc_proto_rawDescData
}

var file_cluster_proto_cluster_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_cluster_proto_cluster_rpc_proto_goTypes = []interface{}{
	(*Node)(nil),        // 0: coms.proto.cluster.Node
	(*NodeList)(nil),    // 1: coms.proto.cluster.NodeList
	(*Void)(nil),        // 2: coms.proto.cluster.Void
	(*Id)(nil),          // 3: coms.proto.cluster.Id
	(*ExchangeReq)(nil), // 4: coms.proto.cluster.ExchangeReq
}
var file_cluster_proto_cluster_rpc_proto_depIdxs = []int32{
	0, // 0: coms.proto.cluster.NodeList.nodes:type_name -> coms.proto.cluster.Node
	1, // 1: coms.proto.cluster.ClusterRpcService.Gossip:input_type -> coms.proto.cluster.NodeList
	2, // 2: coms.proto.cluster.ClusterRpcService.GetId:input_type -> coms.proto.cluster.Void
	4, // 3: coms.proto.cluster.ClusterRpcService.Exchange:input_type -> coms.proto.cluster.ExchangeReq
	2, // 4: coms.proto.cluster.ClusterRpcService.Gossip:output_type -> coms.proto.cluster.Void
	3, // 5: coms.proto.cluster.ClusterRpcService.GetId:output_type -> coms.proto.cluster.Id
	2, // 6: coms.proto.cluster.ClusterRpcService.Exchange:output_type -> coms.proto.cluster.Void
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_cluster_proto_cluster_rpc_proto_init() }
func file_cluster_proto_cluster_rpc_proto_init() {
	if File_cluster_proto_cluster_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cluster_proto_cluster_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cluster_proto_cluster_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cluster_proto_cluster_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Void); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cluster_proto_cluster_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Id); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cluster_proto_cluster_rpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExchangeReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_cluster_proto_cluster_rpc_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cluster_proto_cluster_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cluster_proto_cluster_rpc_proto_goTypes,
		DependencyIndexes: file_cluster_proto_cluster_rpc_proto_depIdxs,
		MessageInfos:      file_cluster_proto_cluster_rpc_proto_msgTypes,
	}.Build()
	File_cluster_proto_cluster_rpc_proto = out.File
	file_cluster_proto_cluster_rpc_proto_rawDesc = nil
	file_cluster_proto_cluster_rpc_proto_goTypes = nil
	file_cluster_proto_cluster_rpc_proto_depIdxs = nil
}
