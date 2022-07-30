package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"simple-napster/dal"
	"simple-napster/entities"
	messages "simple-napster/protos/messages"

	services "simple-napster/protos/services"
)

type NapsterServerListener struct {
	services.UnimplementedNapsterServer
	dal dal.ServerDal

	grpcServer *grpc.Server
	listener   net.Listener
}

type NapsterServerListenerConfig struct {
	Port int
}

func NewNapsterServerListener(
	config *NapsterServerListenerConfig,
	dal dal.ServerDal,
) *NapsterServerListener {
	napster := &NapsterServerListener{dal: dal}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}
	napster.listener = listener

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	services.RegisterNapsterServer(grpcServer, napster)

	napster.grpcServer = grpcServer

	return napster
}

func (sl *NapsterServerListener) Start() {
	if err := sl.grpcServer.Serve(sl.listener); err != nil {
		panic(err)
	}
}

// Join TODO: check if bun correctly sets IDs and see how is the data in the database after a JOIN request
func (sl *NapsterServerListener) Join(ctx context.Context, args *messages.JoinArgs) (*messages.JoinResponse, error) {
	fmt.Printf("Join request received \n peer_ip: %s \n peer_port: %d \n", args.IP, args.Port)

	peer := &entities.Peer{
		Port:   args.Port,
		Active: true,
		IP:     args.IP,
	}
	err := sl.dal.AddPeer(ctx, peer)
	if err != nil {
		fmt.Printf("Error when trying to persist peer: %s", err.Error())
		return nil, err
	}

	files := make([]*entities.File, len(args.Files))
	for i, f := range args.Files {
		files[i] = &entities.File{Name: f}
	}

	err = sl.dal.AddFilesIfNew(ctx, files)
	if err != nil {
		fmt.Printf("Error when adding new peer files: %s", err.Error())
		return nil, err
	}

	err = sl.dal.AddFilesToPeer(ctx, peer, files)
	if err != nil {
		fmt.Printf("Error when associating files with peer: %s", err.Error())
		return nil, err
	}
	return &messages.JoinResponse{Id: peer.ID.String()}, nil
}

func (sl *NapsterServerListener) Leave(ctx context.Context, args *messages.LeaveArgs) (*emptypb.Empty, error) {
	fmt.Printf("Leave request received peer_id: %s\n", args.PeerId)
	id, err := uuid.Parse(args.PeerId)
	if err != nil {
		fmt.Println("Invalid peer id")
		return nil, err
	}

	err = sl.dal.DeletePeerAndFiles(ctx, id)
	if err != nil {
		fmt.Printf("Error when deleting peer files: %s \n", err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (sl *NapsterServerListener) Search(ctx context.Context, args *messages.SearchArgs) (*messages.SearchResponse, error) {
	fmt.Printf("Search request received for file: %s \n", args.Filename)

	peers, err := sl.dal.GetAllPeersWithFile(ctx, args.Filename)
	if err != nil {

		return nil, err
	}

	protos := make([]*messages.Peer, len(peers))
	for i, peer := range peers {
		protos[i] = &messages.Peer{
			Port: peer.Port,
			IP:   peer.IP,
		}
	}

	return &messages.SearchResponse{AvailablePeers: protos}, nil
}

func (sl *NapsterServerListener) Update(ctx context.Context, args *messages.UpdateArgs) (*emptypb.Empty, error) {
	fmt.Printf("Update request received for peer %s \n", args.PeerId)
	id, err := uuid.Parse(args.PeerId)
	if err != nil {
		return nil, err
	}
	peer := &entities.Peer{ID: id}

	err = sl.dal.AddFileToPeeerWithFilename(ctx, peer, args.NewFile)
	if err != nil {
		fmt.Printf("Error when persisting files for peer: %s \n", err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
