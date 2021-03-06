// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: protos/napster_peer_service.proto

package services

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	messages "simple-napster/protos/messages"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NapsterPeerClient is the client API for NapsterPeer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NapsterPeerClient interface {
	IsAlive(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DownloadFile(ctx context.Context, in *messages.DownloadFileArgs, opts ...grpc.CallOption) (NapsterPeer_DownloadFileClient, error)
}

type napsterPeerClient struct {
	cc grpc.ClientConnInterface
}

func NewNapsterPeerClient(cc grpc.ClientConnInterface) NapsterPeerClient {
	return &napsterPeerClient{cc}
}

func (c *napsterPeerClient) IsAlive(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/NapsterPeer/IsAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *napsterPeerClient) DownloadFile(ctx context.Context, in *messages.DownloadFileArgs, opts ...grpc.CallOption) (NapsterPeer_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &NapsterPeer_ServiceDesc.Streams[0], "/NapsterPeer/DownloadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &napsterPeerDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NapsterPeer_DownloadFileClient interface {
	Recv() (*messages.DownloadFileResponse, error)
	grpc.ClientStream
}

type napsterPeerDownloadFileClient struct {
	grpc.ClientStream
}

func (x *napsterPeerDownloadFileClient) Recv() (*messages.DownloadFileResponse, error) {
	m := new(messages.DownloadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NapsterPeerServer is the server API for NapsterPeer service.
// All implementations must embed UnimplementedNapsterPeerServer
// for forward compatibility
type NapsterPeerServer interface {
	IsAlive(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	DownloadFile(*messages.DownloadFileArgs, NapsterPeer_DownloadFileServer) error
	mustEmbedUnimplementedNapsterPeerServer()
}

// UnimplementedNapsterPeerServer must be embedded to have forward compatible implementations.
type UnimplementedNapsterPeerServer struct {
}

func (UnimplementedNapsterPeerServer) IsAlive(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAlive not implemented")
}
func (UnimplementedNapsterPeerServer) DownloadFile(*messages.DownloadFileArgs, NapsterPeer_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedNapsterPeerServer) mustEmbedUnimplementedNapsterPeerServer() {}

// UnsafeNapsterPeerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NapsterPeerServer will
// result in compilation errors.
type UnsafeNapsterPeerServer interface {
	mustEmbedUnimplementedNapsterPeerServer()
}

func RegisterNapsterPeerServer(s grpc.ServiceRegistrar, srv NapsterPeerServer) {
	s.RegisterService(&NapsterPeer_ServiceDesc, srv)
}

func _NapsterPeer_IsAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NapsterPeerServer).IsAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NapsterPeer/IsAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NapsterPeerServer).IsAlive(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _NapsterPeer_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(messages.DownloadFileArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NapsterPeerServer).DownloadFile(m, &napsterPeerDownloadFileServer{stream})
}

type NapsterPeer_DownloadFileServer interface {
	Send(*messages.DownloadFileResponse) error
	grpc.ServerStream
}

type napsterPeerDownloadFileServer struct {
	grpc.ServerStream
}

func (x *napsterPeerDownloadFileServer) Send(m *messages.DownloadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

// NapsterPeer_ServiceDesc is the grpc.ServiceDesc for NapsterPeer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NapsterPeer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NapsterPeer",
	HandlerType: (*NapsterPeerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsAlive",
			Handler:    _NapsterPeer_IsAlive_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadFile",
			Handler:       _NapsterPeer_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protos/napster_peer_service.proto",
}
