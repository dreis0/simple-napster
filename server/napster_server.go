package main

import (
	"simple-napster/dal"
	"sync"
)

type NapsterServer struct {
	server *NapsterServerListener
}

func NewNapsterServer(dal dal.ServerDal, port int) *NapsterServer {
	// TODO: instantiate keep alive client
	return &NapsterServer{
		server: NewNapsterServerListener(&NapsterServerListenerConfig{Port: port}, dal),
	}
}

func (s *NapsterServer) Start() {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	go func() {
		s.server.Start()
		waitGroup.Done()
	}()
}
