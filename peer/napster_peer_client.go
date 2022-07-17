package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"os"
	"simple-napster/entities"
	messages "simple-napster/protos/messages"
	napsterProto "simple-napster/protos/services"
	"strconv"
	"strings"
)

type NapsterPeerClient struct {
	client   napsterProto.NapsterClient
	selfId   string
	selfPort int
	filePath string
}

type NapsterPeerClientConfig struct {
	ServerIp   string
	ServerPort int
	SelfPort   int
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
		selfPort: config.SelfPort,
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
		case "SEARCH":
			c.SearchRequest(ctx, reader)
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

	args := &messages.JoinArgs{IP: "localhost", Port: int32(c.selfPort), Files: filenames}
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

func (c *NapsterPeerClient) SearchRequest(ctx context.Context, reader *bufio.Reader) {
	fmt.Println("Enter filename:")
	filename := readInput(reader)

	args := &messages.SearchArgs{Filename: filename}
	response, err := c.client.Search(ctx, args)
	if err != nil {
		fmt.Printf("Fail to perform SEARCH action: %s \n", err.Error())
		return
	}
	if len(response.AvailablePeers) == 0 {
		fmt.Println("No peers available")
		return
	}

	fmt.Println("Choose a peer by typing the corresponding index: ")
	fmt.Println("0 - Cancel")
	for i, p := range response.AvailablePeers {
		fmt.Printf("%d - %s:%d\n", i+1, p.IP, p.Port)
	}

	peerIdx, err := strconv.Atoi(readInput(reader))
	if err != nil {
		fmt.Println("Invalid input")
		return
	}
	if peerIdx != 0 {
		c.DownloadRequest(ctx, &entities.Peer{
			IP:   response.AvailablePeers[peerIdx-1].IP,
			Port: response.AvailablePeers[peerIdx-1].Port,
		}, filename)
	}
}

func (c *NapsterPeerClient) DownloadRequest(ctx context.Context, peer *entities.Peer, filename string) {
	peerAddress := fmt.Sprintf("%s:%d", peer.IP, peer.Port)
	conn, err := grpc.Dial(peerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Fail to perform DOWNLOAD action: %s \n", err.Error())
		return
	}

	client := napsterProto.NewNapsterPeerClient(conn)
	stream, err := client.DownloadFile(ctx, &messages.DownloadFileArgs{Partition: 1, FileName: filename})
	if err != nil {
		fmt.Printf("Fail to perform DOWNLOAD action: %s \n", err.Error())
		return
	}

	done := make(chan bool)
	fileBytes := []byte{}

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				done <- true //means stream is finished
				return
			}
			if err != nil {
				fmt.Printf("Fail to perform DOWNLOAD action. Download interrupted: %s \n", err.Error())
				done <- true
				break
			}
			fileBytes = append(fileBytes, response.FileBytes...)
		}
	}()

	var d bool
	d = <-done

	// TODO: trigger update request after download
	if d {
		file, err := os.Create(c.filePath + "/" + filename)
		if err != nil {
			fmt.Printf("Fail to perform DOWNLOAD action. Failed to create file: %s \n", err.Error())
			returnw
		}
		defer file.Close()
		file.Write(fileBytes)

		fmt.Println("DONWLOAD_OK")
	} else {
		fmt.Println("DOWNLOAD_FAILED")
	}
}

func readInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return strings.Replace(input, "\n", "", -1)
}
