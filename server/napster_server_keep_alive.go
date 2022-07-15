package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"simple-napster/dal"
	napsterproto "simple-napster/protos/services"
)

type NapsterServerKeepAliveClient struct {
	dal               dal.ServerDal
	failedAttemptsMap map[string]int // TODO: Use this to keep track of failed calls and decide whether to inactivate the peer
	ctx               context.Context
}

func NewKeepAliveClient(dal dal.ServerDal, ctx context.Context) *NapsterServerKeepAliveClient {
	return &NapsterServerKeepAliveClient{dal: dal, ctx: ctx}
}

func (ka *NapsterServerKeepAliveClient) Start() {
	peers, err := ka.dal.GetPeers(ka.ctx)
	if err != nil {
		fmt.Errorf("Failed to get peers", err)
	}

	for _, peer := range peers {
		peerAddress := fmt.Sprintf("%s:%d", peer.IP, peer.Port)
		conn, err := grpc.Dial(peerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}

		client := napsterproto.NewNapsterPeerClient(conn)
		_, _ = client.IsAlive(ka.ctx, &emptypb.Empty{})

	}
}
