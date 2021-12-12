// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	GetChat(ctx context.Context, in *GetChatRequest, opts ...grpc.CallOption) (*Chat, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) GetChat(ctx context.Context, in *GetChatRequest, opts ...grpc.CallOption) (*Chat, error) {
	out := new(Chat)
	err := c.cc.Invoke(ctx, "/api.proc.ChatService/GetChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	GetChat(context.Context, *GetChatRequest) (*Chat, error)
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (*UnimplementedChatServiceServer) GetChat(context.Context, *GetChatRequest) (*Chat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChat not implemented")
}
func (*UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

func RegisterChatServiceServer(s *grpc.Server, srv ChatServiceServer) {
	s.RegisterService(&_ChatService_serviceDesc, srv)
}

func _ChatService_GetChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.proc.ChatService/GetChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetChat(ctx, req.(*GetChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChatService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.proc.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChat",
			Handler:    _ChatService_GetChat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}
