package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"simple-napster/dal"
	"simple-napster/protos/messages"
	services "simple-napster/protos/services"
)

type NapsterPeerServer struct {
	services.UnimplementedNapsterPeerServer

	dal      dal.ClientDal
	filePath string

	listener   net.Listener
	grpcServer *grpc.Server
}

type NapsterPeerServerConfig struct {
	Port     int
	FilePath string
}

func NewNapsterPeerServer(config *NapsterPeerServerConfig, dal dal.ClientDal) *NapsterPeerServer {
	server := &NapsterPeerServer{dal: dal, filePath: config.FilePath}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}
	server.listener = listener

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	services.RegisterNapsterPeerServer(grpcServer, server)

	server.grpcServer = grpcServer

	return server
}

func (ps *NapsterPeerServer) Start() {
	if err := ps.grpcServer.Serve(ps.listener); err != nil {
		panic(err)
	}
}

func (peer *NapsterPeerServer) IsAlive(ctx context.Context, args *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (peer *NapsterPeerServer) DownloadFile(args *messages.DownloadFileArgs, server services.NapsterPeer_DownloadFileServer) error {
	if !shouldAcceptRequest() {
		return errors.New("DOWNLOAD_NEGADO")
	}
	files, err := ioutil.ReadDir(peer.filePath)
	if err != nil {
		return err
	}

	fileIdx := -1
	for i, fi := range files {
		if fi.Name() == args.FileName {
			fileIdx = i
		}
	}

	if fileIdx == -1 {
		return errors.New("file not found")
	}

	fileInfo := files[fileIdx]
	file, err := os.Open(peer.filePath + "/" + fileInfo.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	size := fileInfo.Size()
	ammountOfSends := size / 10000

	for i := 0; i <= int(ammountOfSends); i++ {
		bytes := make([]byte, size/ammountOfSends)
		if err != nil {
			return err
		}
		_, err := file.Read(bytes)
		if err != nil && err != io.EOF {
			return err
		}

		server.Send(&messages.DownloadFileResponse{
			FileBytes:       bytes,
			Partition:       int32(i),
			TotalPartitions: int32(ammountOfSends),
		})
	}

	return nil
}

func shouldAcceptRequest() bool {
	return rand.Intn(1) == 0
}
