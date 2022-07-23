package main

import (
	"github.com/joho/godotenv"
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

	portNum, err := strconv.Atoi(portStr)
	if err != nil {
		panic("invalid port number")
	}

	files_path, err := utils.GetArgument(args, "files-path")
	if err != nil {
		panic("no path provided to read and save files")
	}

	serverIp, err := utils.GetArgument(args, "server-ip")
	if err != nil {
		panic("no server IP provided")
	}

	serverPortStr, err := utils.GetArgument(args, "server-port")
	if err != nil {
		panic("no server port provided")
	}

	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		panic("invalid server port provided")
	}

	dal := createDal()
	peer := NewNapsterPeer(
		&NapsterPeerConfig{
			ServerIp:   serverIp,
			ServerPort: serverPort,
			SelfPort:   portNum,
			FilePath:   files_path,
		},
		dal,
	)

	peer.Start()
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
