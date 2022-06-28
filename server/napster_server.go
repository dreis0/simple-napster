package main

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"simple-napster/dal"
	"simple-napster/entities"
	messages "simple-napster/protos/messages"
	services "simple-napster/protos/services"
)

type NapsterServer struct {
	services.UnimplementedNapsterServer
	dal dal.Dal
}

func NewNapsterService(dal dal.Dal) *NapsterServer {
	return &NapsterServer{dal: dal}
}

func (s *NapsterServer) Join(ctx context.Context, args *messages.JoinArgs) (*emptypb.Empty, error) {
	peer := &entities.Peer{
		Port:   args.Port,
		Active: true,
		IP:     args.IP,
	}
	err := s.dal.AddPeer(ctx, peer)

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
