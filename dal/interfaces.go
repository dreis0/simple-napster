package dal

import (
	"context"
	"github.com/google/uuid"
	"simple-napster/entities"
)

type ServerDal interface {
	AddPeer(ctx context.Context, peer *entities.Peer) error
	GetPeers(ctx context.Context) ([]entities.Peer, error)
	InativatePeer(ctx context.Context, peerId uuid.UUID) error
}

type ClientDal interface {
}
