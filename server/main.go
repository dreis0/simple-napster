package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"simple-napster/dal"
	"simple-napster/utils"
	"strconv"
)

func main() {
	godotenv.Load(".env")
	
	args := os.Args
	portStr, err := utils.GetArgument(args, "port")
	if err != nil {
		panic("no port provided")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("invalid port provided")
	}
	log.Printf("port: %d", port)

	server := NewNapsterServer(dal.NewDalFromEnv(), port)
	server.Start()
}
