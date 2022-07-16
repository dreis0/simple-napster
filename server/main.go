package main

import (
	"log"
	"os"
	"simple-napster/dal"
	"simple-napster/utils"
	"strconv"
)

func main() {

	args := os.Args
	portStr, err := utils.GetArgument(args, "port")
	if err != nil {
		panic("no port provided")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("invalid port provided")
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
