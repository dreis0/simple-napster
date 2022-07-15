package main

import (
	"context"
	"simple-napster/dal"
	"sync"
)

type NapsterServer struct {
	server    *NapsterServerListener
	keepAlive *NapsterServerKeepAliveClient
}

func NewNapsterServer(dal dal.ServerDal, port int) *NapsterServer {
	return &NapsterServer{
		server:    NewNapsterServerListener(&NapsterServerListenerConfig{Port: port}, dal),
		keepAlive: NewKeepAliveClient(dal, context.Background()),
	}
}

func (s *NapsterServer) Start() {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	go func() {
		s.server.Start()
		waitGroup.Done()
	}()

	go func() {
		s.keepAlive.Start()
		waitGroup.Done()
	}()

	waitGroup.Wait()
}
