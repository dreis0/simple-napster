package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"simple-napster/dal"
	services "simple-napster/protos/services"
)

type NapsterPeerServer struct {
	services.UnimplementedNapsterPeerServer

	dal dal.ClientDal

	listener   net.Listener
	grpcServer *grpc.Server
}

func NewNapsterPeerServer(port int, dal dal.ClientDal) *NapsterPeerServer {
	server := &NapsterPeerServer{dal: dal}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	server.listener = listener

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	services.RegisterNapsterPeerServer(grpcServer, server)

	return server
}

func (ps *NapsterPeerServer) Start() {
	if err := ps.grpcServer.Serve(ps.listener); err != nil {
		panic(err)
	}
}

func (peer *NapsterPeerServer) IsAlive(ctx context.Context, args *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("Is alive received")
	return &emptypb.Empty{}, nil
}
