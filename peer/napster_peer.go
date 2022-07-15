package main

import (
	"fmt"
	"simple-napster/dal"
	"sync"
)

type NapsterPeerConfig struct {
	SelfPort   int
	ServerIp   string
	ServerPort int
}

type NapsterPeer struct {
	client *NapsterPeerClient
	server *NapsterPeerServer
}

func NewNapsterPeer(config *NapsterPeerConfig, dal dal.ClientDal) *NapsterPeer {
	return &NapsterPeer{
		client: NewPeerClient(&NapsterPeerClientConfig{ServerIp: config.ServerIp, ServerPort: config.ServerPort}),
		server: NewNapsterPeerServer(config.SelfPort, dal),
	}
}

func (peer *NapsterPeer) Start() {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		peer.client.Start()
	}()
	go func() {
		defer waitGroup.Done()
		peer.server.Start()
	}()

	waitGroup.Wait()

	fmt.Println("The peer has died")
}
