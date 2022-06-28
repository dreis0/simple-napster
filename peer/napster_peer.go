package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strings"

	messages "simple-napster/protos/messages"
	napsterProto "simple-napster/protos/services"
)

type Config struct {
	SelfPort   int
	ServerIp   string
	ServerPort int
}

type NapsterPeer struct {
	config *Config
}

func NewNapsterPeer(config *Config) *NapsterPeer {
	return &NapsterPeer{
		config: config,
	}
}

func (peer *NapsterPeer) RunClient() {
	serverAddress := fmt.Sprintf("%s:%d", peer.config.ServerIp, peer.config.ServerPort)
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	client := napsterProto.NewNapsterClient(conn)
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for true {
		fmt.Printf("Type your command")
		input := readInput(reader)
		switch input {
		case "JOIN":
			args := &messages.JoinArgs{IP: "localhost", Port: 3000, Files: []string{}}
			_, err := client.Join(ctx, args)

			if err != nil {
				fmt.Printf("Fail to perform JOIN action: %s", err.Error())
			} else {
				fmt.Print("JOIN_OK")
			}
		}
	}
}

func (peer *NapsterPeer) RunServer() {

}

func readInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return strings.Replace(input, "\n", "", -1)
}
