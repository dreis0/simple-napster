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

type KeepAliveClient struct {
	dal               dal.ServerDal
	failedAttemptsMap map[string]int // TODO: Use this to keep track of failed calls and decide whether to inactivate the peer
}

func NewKeepAliveClient(dal dal.ServerDal) *KeepAliveClient {
	return &KeepAliveClient{dal: dal}
}

func (keepAlive *KeepAliveClient) RunKeepAliveForAllPeers(ctx context.Context) {
	peers, err := keepAlive.dal.GetPeers(ctx)
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
		_, _ = client.IsAlive(ctx, &emptypb.Empty{})

	}
}
