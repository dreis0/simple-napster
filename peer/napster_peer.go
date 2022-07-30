package main

import (
	"fmt"
	"sync"
)

type NapsterPeerConfig struct {
	SelfPort   int
	ServerIp   string
	ServerPort int
	FilePath   string
}

type NapsterPeer struct {
	client *NapsterPeerClient
	server *NapsterPeerServer
}

func NewNapsterPeer(config *NapsterPeerConfig) *NapsterPeer {
	return &NapsterPeer{
		client: NewPeerClient(&NapsterPeerClientConfig{SelfPort: config.SelfPort, ServerIp: config.ServerIp, ServerPort: config.ServerPort, FilePath: config.FilePath}),
		server: NewNapsterPeerServer(&NapsterPeerServerConfig{Port: config.SelfPort, FilePath: config.FilePath}),
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
