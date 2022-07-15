package main

import (
	"os"
	"simple-napster/dal"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		panic("no port provided")
	}
	port := args[1]
	portNum, _ := strconv.Atoi(port)

	dal := createDal()
	peer := NewNapsterPeer(&Config{ServerPort: 10098, SelfPort: portNum, ServerIp: "localhost"}, dal)

	peer.Run()
}

func createDal() dal.ClientDal {
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