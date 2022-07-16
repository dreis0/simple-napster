package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"os"
	messages "simple-napster/protos/messages"
	napsterProto "simple-napster/protos/services"
	"strings"
)

type NapsterPeerClient struct {
	client   napsterProto.NapsterClient
	selfId   string
	filePath string
}

type NapsterPeerClientConfig struct {
	ServerIp   string
	ServerPort int
	FilePath   string
}

func NewPeerClient(config *NapsterPeerClientConfig) *NapsterPeerClient {
	serverAddress := fmt.Sprintf("%s:%d", config.ServerIp, config.ServerPort)
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return &NapsterPeerClient{
		client:   napsterProto.NewNapsterClient(conn),
		filePath: config.FilePath,
	}
}

func (c *NapsterPeerClient) Start() {
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for true {
		fmt.Println("Type your command")
		input := readInput(reader)

		switch input {
		case "JOIN":
			c.JoinRequest(ctx)
			break
		case "LEAVE":
			c.LeaveRequest(ctx)
			break
		}
	}
}

func (c *NapsterPeerClient) JoinRequest(ctx context.Context) {
	files, err := ioutil.ReadDir(c.filePath)
	if err != nil {
		fmt.Printf("Fail to perform JOIN action. Cannot read files: %s \n", err.Error())
	}

	filenames := make([]string, len(files))
	for i, fi := range files {
		filenames[i] = fi.Name()
	}

	args := &messages.JoinArgs{IP: "localhost", Port: 3000, Files: filenames}
	response, err := c.client.Join(ctx, args)

	if err != nil {
		fmt.Printf("Fail to perform JOIN action: %s \n", err.Error())
	} else {
		c.selfId = response.Id
		fmt.Println("JOIN_OK")
	}
}

func (c *NapsterPeerClient) LeaveRequest(ctx context.Context) {
	args := &messages.LeaveArgs{PeerId: c.selfId}

	_, err := c.client.Leave(ctx, args)
	if err != nil {
		fmt.Printf("Fail to perform LEAVE action: %s \n", err.Error())
	} else {
		fmt.Println("LEAVE_OK")
	}
}

func readInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return strings.Replace(input, "\n", "", -1)
}
