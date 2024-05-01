// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: pubsub.proto

package pubsubPB

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

const (
	PubSubServer_Publish_FullMethodName   = "/pubsub.PubSubServer/Publish"
	PubSubServer_Subscribe_FullMethodName = "/pubsub.PubSubServer/Subscribe"
)

// PubSubServerClient is the client API for PubSubServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PubSubServerClient interface {
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (PubSubServer_SubscribeClient, error)
}

type pubSubServerClient struct {
	cc grpc.ClientConnInterface
}

func NewPubSubServerClient(cc grpc.ClientConnInterface) PubSubServerClient {
	return &pubSubServerClient{cc}
}

func (c *pubSubServerClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, PubSubServer_Publish_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServerClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (PubSubServer_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &PubSubServer_ServiceDesc.Streams[0], PubSubServer_Subscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pubSubServerSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PubSubServer_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type pubSubServerSubscribeClient struct {
	grpc.ClientStream
}

func (x *pubSubServerSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PubSubServerServer is the server API for PubSubServer service.
// All implementations must embed UnimplementedPubSubServerServer
// for forward compatibility
type PubSubServerServer interface {
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	Subscribe(*SubscribeRequest, PubSubServer_SubscribeServer) error
	mustEmbedUnimplementedPubSubServerServer()
}

// UnimplementedPubSubServerServer must be embedded to have forward compatible implementations.
type UnimplementedPubSubServerServer struct {
}

func (UnimplementedPubSubServerServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedPubSubServerServer) Subscribe(*SubscribeRequest, PubSubServer_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedPubSubServerServer) mustEmbedUnimplementedPubSubServerServer() {}

// UnsafePubSubServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PubSubServerServer will
// result in compilation errors.
type UnsafePubSubServerServer interface {
	mustEmbedUnimplementedPubSubServerServer()
}

func RegisterPubSubServerServer(s grpc.ServiceRegistrar, srv PubSubServerServer) {
	s.RegisterService(&PubSubServer_ServiceDesc, srv)
}

func _PubSubServer_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServerServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubServer_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServerServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubServer_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PubSubServerServer).Subscribe(m, &pubSubServerSubscribeServer{stream})
}

type PubSubServer_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type pubSubServerSubscribeServer struct {
	grpc.ServerStream
}

func (x *pubSubServerSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

// PubSubServer_ServiceDesc is the grpc.ServiceDesc for PubSubServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PubSubServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pubsub.PubSubServer",
	HandlerType: (*PubSubServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _PubSubServer_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _PubSubServer_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pubsub.proto",
}
