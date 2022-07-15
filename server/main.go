package main

import (
	"context"
	"log"
	"os"
	"simple-napster/dal"
	"strconv"
	"time"
)

func main() {

	args := os.Args
	if len(args) <= 1 {
		panic("no port provided")
	}

	port, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}
	log.Printf("port: %s", port)

	dal := createDal()

	server := NewNapsterServer(dal, port)
	server.Start()
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
