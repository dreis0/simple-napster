package main

import (
	"log"
	"os"
	"simple-napster/dal"
	"strconv"
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
