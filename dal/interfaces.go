package dal

import (
	"context"
	"simple-napster/entities"
)

type ServerDal interface {
	AddPeer(ctx context.Context, peer *entities.Peer) error
	GetPeers(ctx context.Context) ([]*entities.Peer, error)
	UpdatePeer(ctx context.Context, peer *entities.Peer) error
}

type ClientDal interface {
}
