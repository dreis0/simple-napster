package dal

import (
	"context"
	"github.com/google/uuid"
	"simple-napster/entities"
)

type ServerDal interface {
	AddPeer(ctx context.Context, peer *entities.Peer) error
	GetPeers(ctx context.Context) ([]*entities.Peer, error)
	GetPeerById(ctx context.Context, id uuid.UUID) (*entities.Peer, error)
	DeletePeerAndFiles(ctx context.Context, id uuid.UUID) error
	UpdatePeer(ctx context.Context, peer *entities.Peer) error
	AddFilesIfNew(ctx context.Context, files []*entities.File) error
	AddFilesToPeer(ctx context.Context, peer *entities.Peer, files []*entities.File) error
	GetAllPeersWithFile(ctx context.Context, filename string) ([]*entities.Peer, error)
}

type ClientDal interface {
}
