package main

import (
	"github.com/joho/godotenv"
	"os"
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

	peer := NewNapsterPeer(
		&NapsterPeerConfig{
			ServerIp:   serverIp,
			ServerPort: serverPort,
			SelfPort:   portNum,
			FilePath:   files_path,
		},
	)

	peer.Start()
}
