package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	services "simple-napster/protos/services"
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

	services.RegisterNapsterServer(s, NewNapsterService())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("listening on port %s", args[1])
}
