package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	messages "simple-napster/protos/messages"
	services "simple-napster/protos/services"
)

type NapsterService struct {
	services.UnimplementedNapsterServer
}

func NewNapsterService() *NapsterService {
	return &NapsterService{}
}

func (s *NapsterService) Join(ctx context.Context, args *messages.JoinArgs) (*emptypb.Empty, error) {
	fmt.Println("JOIN REQUEST")
	return &emptypb.Empty{}, nil
}
