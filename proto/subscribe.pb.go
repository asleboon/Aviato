// Code generated by protoc-gen-go. DO NOT EDIT.
// source: subscribe.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// message keyword will generate a struct in golang
// Changed to uint 64 and compiled
type SubscribeMessage struct {
	RefreshRate          uint64   `protobuf:"varint,1,opt,name=refreshRate,proto3" json:"refreshRate,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeMessage) Reset()         { *m = SubscribeMessage{} }
func (m *SubscribeMessage) String() string { return proto.CompactTextString(m) }
func (*SubscribeMessage) ProtoMessage()    {}
func (*SubscribeMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_38d2980c9543da44, []int{0}
}

func (m *SubscribeMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeMessage.Unmarshal(m, b)
}
func (m *SubscribeMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeMessage.Marshal(b, m, deterministic)
}
func (m *SubscribeMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeMessage.Merge(m, src)
}
func (m *SubscribeMessage) XXX_Size() int {
	return xxx_messageInfo_SubscribeMessage.Size(m)
}
func (m *SubscribeMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeMessage proto.InternalMessageInfo

func (m *SubscribeMessage) GetRefreshRate() uint64 {
	if m != nil {
		return m.RefreshRate
	}
	return 0
}

type NotificationMessage struct {
	Top10                string   `protobuf:"bytes,1,opt,name=top10,proto3" json:"top10,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotificationMessage) Reset()         { *m = NotificationMessage{} }
func (m *NotificationMessage) String() string { return proto.CompactTextString(m) }
func (*NotificationMessage) ProtoMessage()    {}
func (*NotificationMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_38d2980c9543da44, []int{1}
}

func (m *NotificationMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotificationMessage.Unmarshal(m, b)
}
func (m *NotificationMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotificationMessage.Marshal(b, m, deterministic)
}
func (m *NotificationMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotificationMessage.Merge(m, src)
}
func (m *NotificationMessage) XXX_Size() int {
	return xxx_messageInfo_NotificationMessage.Size(m)
}
func (m *NotificationMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_NotificationMessage.DiscardUnknown(m)
}

var xxx_messageInfo_NotificationMessage proto.InternalMessageInfo

func (m *NotificationMessage) GetTop10() string {
	if m != nil {
		return m.Top10
	}
	return ""
}

func init() {
	proto.RegisterType((*SubscribeMessage)(nil), "proto.SubscribeMessage")
	proto.RegisterType((*NotificationMessage)(nil), "proto.NotificationMessage")
}

func init() { proto.RegisterFile("subscribe.proto", fileDescriptor_38d2980c9543da44) }

var fileDescriptor_38d2980c9543da44 = []byte{
	// 148 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x2e, 0x4d, 0x2a,
	0x4e, 0x2e, 0xca, 0x4c, 0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a,
	0xea, 0x5c, 0x02, 0xc1, 0x30, 0x19, 0xdf, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0x21, 0x61, 0x2e,
	0xee, 0xa2, 0xd4, 0xb4, 0xa2, 0xd4, 0xe2, 0x8c, 0xa0, 0xc4, 0x92, 0x54, 0x09, 0x46, 0x05, 0x46,
	0x0d, 0x16, 0x25, 0x15, 0x2e, 0x61, 0xbf, 0xfc, 0x92, 0xcc, 0xb4, 0xcc, 0xe4, 0xc4, 0x92, 0xcc,
	0xfc, 0x3c, 0x98, 0x5a, 0x5e, 0x2e, 0xd6, 0x92, 0xfc, 0x02, 0x43, 0x03, 0xb0, 0x2a, 0x4e, 0xa3,
	0x30, 0x2e, 0x1e, 0xa8, 0x71, 0x05, 0x20, 0x55, 0x42, 0x6e, 0x5c, 0x9c, 0x70, 0xe3, 0x85, 0xc4,
	0x21, 0x56, 0xeb, 0xa1, 0x5b, 0x28, 0x25, 0x05, 0x95, 0xc0, 0x62, 0x81, 0x12, 0x83, 0x06, 0xa3,
	0x01, 0x63, 0x12, 0x1b, 0x58, 0x81, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xca, 0xa7, 0x68, 0x12,
	0xc7, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SubscriptionClient is the client API for Subscription service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SubscriptionClient interface {
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (Subscription_SubscribeClient, error)
}

type subscriptionClient struct {
	cc *grpc.ClientConn
}

func NewSubscriptionClient(cc *grpc.ClientConn) SubscriptionClient {
	return &subscriptionClient{cc}
}

func (c *subscriptionClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (Subscription_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Subscription_serviceDesc.Streams[0], "/proto.Subscription/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &subscriptionSubscribeClient{stream}
	return x, nil
}

type Subscription_SubscribeClient interface {
	Send(*SubscribeMessage) error
	Recv() (*NotificationMessage, error)
	grpc.ClientStream
}

type subscriptionSubscribeClient struct {
	grpc.ClientStream
}

func (x *subscriptionSubscribeClient) Send(m *SubscribeMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *subscriptionSubscribeClient) Recv() (*NotificationMessage, error) {
	m := new(NotificationMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SubscriptionServer is the server API for Subscription service.
type SubscriptionServer interface {
	Subscribe(Subscription_SubscribeServer) error
}

func RegisterSubscriptionServer(s *grpc.Server, srv SubscriptionServer) {
	s.RegisterService(&_Subscription_serviceDesc, srv)
}

func _Subscription_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SubscriptionServer).Subscribe(&subscriptionSubscribeServer{stream})
}

type Subscription_SubscribeServer interface {
	Send(*NotificationMessage) error
	Recv() (*SubscribeMessage, error)
	grpc.ServerStream
}

type subscriptionSubscribeServer struct {
	grpc.ServerStream
}

func (x *subscriptionSubscribeServer) Send(m *NotificationMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *subscriptionSubscribeServer) Recv() (*SubscribeMessage, error) {
	m := new(SubscribeMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Subscription_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Subscription",
	HandlerType: (*SubscriptionServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Subscription_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "subscribe.proto",
}
