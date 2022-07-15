package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"simple-napster/dal"
	services "simple-napster/protos/services"
)

type NapsterPeerServer struct {
	services.UnimplementedNapsterPeerServer
	dal dal.ClientDal
}

func NewNapsterPeerServer(dal dal.ClientDal) *NapsterPeerServer {
	return &NapsterPeerServer{dal: dal}
}

func (peer *NapsterPeerServer) IsAlive(ctx context.Context, args *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("Is alive received")
	return &emptypb.Empty{}, nil
}
