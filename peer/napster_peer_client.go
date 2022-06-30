package main

import (
	"context"
	"fmt"
	messages "simple-napster/protos/messages"
	napsterProto "simple-napster/protos/services"
)

type PeerClient struct {
	client napsterProto.NapsterClient
}

func NewPeerClient(client napsterProto.NapsterClient) *PeerClient {
	return &PeerClient{
		client: client,
	}
}

func (c *PeerClient) DoJoinRequest(ctx context.Context) {
	args := &messages.JoinArgs{IP: "localhost", Port: 3000, Files: []string{}}
	_, err := c.client.Join(ctx, args)

	if err != nil {
		fmt.Printf("Fail to perform JOIN action: %s", err.Error())
	} else {
		fmt.Print("JOIN_OK")
	}
}
