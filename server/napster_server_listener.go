package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"simple-napster/dal"
	"simple-napster/entities"
	messages "simple-napster/protos/messages"

	services "simple-napster/protos/services"
)

type NapsterServerListener struct {
	services.UnimplementedNapsterServer
	dal dal.ServerDal

	grpcServer *grpc.Server
	listener   net.Listener
}

type NapsterServerListenerConfig struct {
	Port int
}

func NewNapsterServerListener(
	config *NapsterServerListenerConfig,
	dal dal.ServerDal,
) *NapsterServerListener {
	napster := &NapsterServerListener{dal: dal}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}
	napster.listener = listener

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	services.RegisterNapsterServer(grpcServer, napster)

	napster.grpcServer = grpcServer

	return napster
}

func (sl *NapsterServerListener) Start() {
	if err := sl.grpcServer.Serve(sl.listener); err != nil {
		panic(err)
	}
}

func (sl *NapsterServerListener) Join(ctx context.Context, args *messages.JoinArgs) (*emptypb.Empty, error) {
	peer := &entities.Peer{
		Port:   args.Port,
		Active: true,
		IP:     args.IP,
	}
	err := sl.dal.AddPeer(ctx, peer)

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
