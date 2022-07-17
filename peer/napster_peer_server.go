package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"io/ioutil"
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

	bytes := make([]byte, fileInfo.Size())
	file.Read(bytes)
	server.Send(&messages.DownloadFileResponse{FileBytes: bytes})

	return nil
}
