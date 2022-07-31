package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"os"
	"simple-napster/entities"
	messages "simple-napster/protos/messages"
	napsterProto "simple-napster/protos/services"
	"simple-napster/utils"
	"strings"
	"time"
)

var (
	ERR_NO_PEERS_AVAILABLE = errors.New("no peers available for specified file")
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

	fmt.Println("Available commands: ")
	fmt.Println("JOIN - join network ")
	fmt.Println("LEAVE - leave network ")
	fmt.Println("DOWNLOAD - attempt download of file")

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
		case "DOWNLOAD":
			c.attemptDownload(ctx, reader)
		}
	}

}

func (c *NapsterPeerClient) JoinRequest(ctx context.Context) {
	if c.selfId != "" {
		fmt.Println("Already in network.")
		return
	}

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
	if !c.isInNetwork() {
		fmt.Println("Peer is not in any network")
		return
	}

	args := &messages.LeaveArgs{PeerId: c.selfId}

	_, err := c.client.Leave(ctx, args)
	if err != nil {
		fmt.Printf("Fail to perform LEAVE action: %s \n", err.Error())
	} else {
		fmt.Println("LEAVE_OK")
	}
}

func (c *NapsterPeerClient) SearchRequest(ctx context.Context, reader *bufio.Reader) ([]*messages.Peer, string, error) {
	fmt.Println("Enter filename:")
	filename := readInput(reader)

	args := &messages.SearchArgs{Filename: filename}
	response, err := c.client.Search(ctx, args)
	if err != nil {
		return nil, "", errors.New(fmt.Sprintf("Fail to perform SEARCH action: %s \n", err.Error()))
	}
	if len(response.AvailablePeers) == 0 {
		return nil, "", ERR_NO_PEERS_AVAILABLE
	}

	return response.AvailablePeers, filename, nil
}

func (c *NapsterPeerClient) DownloadRequest(ctx context.Context, peer *entities.Peer, filename string) error {
	client, err := createPeerStreamClient(peer)

	if err != nil {
		return errors.New(fmt.Sprintf("Fail to perform DOWNLOAD action. Failed to create client: %s \n", err.Error()))
	}

	stream, err := client.DownloadFile(ctx, &messages.DownloadFileArgs{FileName: filename})
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to perform DOWNLOAD action: %s \n", err.Error()))
	}

	done := make(chan bool)
	var downloadErr error
	fileBytes := []byte{}

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				done <- true //means stream is finished
				return
			}
			if err != nil {
				downloadErr = err
				done <- true
				break
			}
			fileBytes = append(fileBytes, response.FileBytes...)
			fmt.Printf("Downloading %d/%d... \n", response.Partition, response.TotalPartitions)
		}
	}()

	<-done
	if downloadErr != nil {
		return errors.New(fmt.Sprintf("Fail to perform DOWNLOAD action: %s \n", downloadErr.Error()))
	}

	// TODO: trigger update request after download
	file, err := os.Create(c.filePath + "/" + filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to perform DOWNLOAD action. Failed to create file: %s \n", err.Error()))
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to perform DOWNLOAD action. Failed to create file: %s \n", err.Error()))
	}

	return nil
}

func (c *NapsterPeerClient) UpdateRequest(ctx context.Context, file string) {
	args := &messages.UpdateArgs{PeerId: c.selfId, NewFile: file}
	c.client.Update(ctx, args)
}

func (c *NapsterPeerClient) attemptDownload(ctx context.Context, reader *bufio.Reader) {
	if !c.isInNetwork() {
		fmt.Println("Peer is not in any network")
		return
	}

	peers, filename, err := c.SearchRequest(ctx, reader)
	if err == ERR_NO_PEERS_AVAILABLE {
		fmt.Println("No peers available")
		return
	}
	if err != nil {
		fmt.Printf("Fail to perform DOWNLOAD %s \n", err.Error())
		return
	}

	maxAttempts := len(peers) * 3
	unusedPeers := make([]*messages.Peer, len(peers))
	copy(unusedPeers, peers)

	for i := 0; i < maxAttempts; i++ {
		peer := &messages.Peer{}
		peer, unusedPeers = selectPeer(unusedPeers)

		fmt.Printf("Attempting to download from %s:%d\n", peer.IP, peer.Port)

		err = c.DownloadRequest(ctx, &entities.Peer{
			IP:   peer.IP,
			Port: peer.Port,
		}, filename)

		if err == nil {
			fmt.Println("DOWNLOAD_OK")
			c.UpdateRequest(ctx, filename)
			break
		}
		if strings.Contains(err.Error(), utils.DOWNLOAD_NEGADO.Error()) {
			fmt.Println("Download denied")
		} else {
			fmt.Printf("Failed to perform DOWNLOAD: %s \n", err.Error())
		}

		if len(unusedPeers) == 0 {
			if i == maxAttempts-1 || len(peers) == 0 {
				fmt.Println("Tried every peer available three times. Unable to perform DOWNLOAD")
			} else {
				unusedPeers = append(unusedPeers, peers...)
				fmt.Println("Failed to download for every peer. Will retry in 30 seconds")
				time.Sleep(30 * time.Second)
			}
		}
	}
}

func createPeerStreamClient(peer *entities.Peer) (napsterProto.NapsterPeerClient, error) {
	peerAddress := fmt.Sprintf("%s:%d", peer.IP, peer.Port)
	conn, err := grpc.Dial(peerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Fail to perform DOWNLOAD action: %s \n", err.Error())
		return nil, err
	}

	return napsterProto.NewNapsterPeerClient(conn), nil
}

func selectPeer(peers []*messages.Peer) (*messages.Peer, []*messages.Peer) {
	idx := utils.RandomInt(len(peers) - 1)
	peer := peers[idx]

	withoutSelected := []*messages.Peer{}
	for i, p := range peers {
		if i != idx {
			withoutSelected = append(withoutSelected, p)
		}
	}

	return peer, withoutSelected
}

func readInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return strings.Replace(input, "\n", "", -1)
}

func (c *NapsterPeerClient) isInNetwork() bool {
	return c.selfId != ""
}
