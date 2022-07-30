package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"simple-napster/dal"
	napsterproto "simple-napster/protos/services"
	"time"
)

type NapsterServerKeepAliveClient struct {
	dal               dal.ServerDal
	failedAttemptsMap map[string]int // TODO: Use this to keep track of failed calls and decide whether to inactivate the peer
	ctx               context.Context
}

func NewKeepAliveClient(dal dal.ServerDal, ctx context.Context) *NapsterServerKeepAliveClient {
	return &NapsterServerKeepAliveClient{dal: dal, ctx: ctx, failedAttemptsMap: make(map[string]int)}
}

func (ka *NapsterServerKeepAliveClient) Start() {
	for true {
		peers, err := ka.dal.GetPeers(ka.ctx)
		if err != nil {
			fmt.Errorf("Failed to get peers", err)
		}

		for _, peer := range peers {
			if !peer.Active {
				continue
			}

			client, err := createGrpcClient(peer.IP, peer.Port)
			if err != nil {
				continue
			}

			_, err = client.IsAlive(ka.ctx, &emptypb.Empty{})
			if err == nil {
				continue
			} else {
				ka.failedAttemptsMap[peer.ID.String()] += 1
				if ka.failedAttemptsMap[peer.ID.String()] >= 3 {
					ka.failedAttemptsMap[peer.ID.String()] = 0
					peer.Active = false
					err = ka.dal.DeletePeerAndFiles(ka.ctx, peer.ID)
					if err != nil {
						fmt.Printf("Fail to inactivate peer %s error: %s \n", peer.ID.String(), err.Error())
					}
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func createGrpcClient(ip string, port int32) (napsterproto.NapsterPeerClient, error) {
	peerAddress := fmt.Sprintf("%s:%d", ip, port)
	conn, err := grpc.Dial(peerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := napsterproto.NewNapsterPeerClient(conn)
	return client, nil
}
