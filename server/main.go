package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"simple-napster/dal"
	services "simple-napster/protos/services"
	"sync"
	"time"
)

func main() {

	args := os.Args
	if len(args) <= 1 {
		panic("no port provided")
	}

	port := args[1]
	log.Printf("port: %s", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	dal := createDal()
	napsterServer := NewNapsterService(dal)
	services.RegisterNapsterServer(s, napsterServer)

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	go func() {
		runKeepAlive(dal)
		defer waitGroup.Done()
	}()

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)

		}
		defer waitGroup.Done()
	}()

	waitGroup.Wait()
}

func createDal() dal.ServerDal {
	config := &dal.Config{
		Hostname: "localhost",
		Port:     5432,
		Password: "postgres",
		Username: "postgres",
		Database: "napster",
	}
	dal := dal.NewDal(config)
	return dal
}

func runKeepAlive(serverDal dal.ServerDal) {
	for true {
		kac := NewKeepAliveClient(serverDal)
		ctx := context.Background()
		kac.RunKeepAliveForAllPeers(ctx)
		time.Sleep(30 * time.Second)
	}
}
