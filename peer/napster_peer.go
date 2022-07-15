package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"simple-napster/dal"
	napsterProto "simple-napster/protos/services"
	services "simple-napster/protos/services"
	"strings"
	"sync"
)

type Config struct {
	SelfPort   int
	ServerIp   string
	ServerPort int
}

type NapsterPeer struct {
	config *Config
	dal    dal.ClientDal
}

func NewNapsterPeer(config *Config, dal dal.ClientDal) *NapsterPeer {
	return &NapsterPeer{
		config: config,
		dal:    dal,
	}
}

func (peer *NapsterPeer) Run() {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		peer.runServer()
	}()
	go func() {
		defer waitGroup.Done()
		peer.runClient()
	}()

	waitGroup.Wait()

	fmt.Println("The peer has died")
}

func (peer *NapsterPeer) runClient() {
	serverAddress := fmt.Sprintf("%s:%d", peer.config.ServerIp, peer.config.ServerPort)
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	client := NewPeerClient(napsterProto.NewNapsterClient(conn))
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for true {
		fmt.Println("Type your command")
		input := readInput(reader)

		switch input {
		case "JOIN":
			client.JoinRequest(ctx)
		}
	}
}

func (peer *NapsterPeer) runServer() {
	s := grpc.NewServer()
	reflection.Register(s)
	peerServer := NewNapsterPeerServer(peer.dal)
	services.RegisterNapsterPeerServer(s, peerServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", peer.config.SelfPort))
	if err != nil {
		panic(err)
	}

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func readInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return strings.Replace(input, "\n", "", -1)
}
