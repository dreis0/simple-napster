package main

import (
	"os"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		panic("no port provided")
	}
	port := args[1]
	portNum, _ := strconv.Atoi(port)

	peer := NewNapsterPeer(&Config{ServerPort: 10098, SelfPort: portNum, ServerIp: "localhost"})
	peer.RunClient()
}
