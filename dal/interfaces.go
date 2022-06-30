package dal

import (
	"context"
	"simple-napster/entities"
)

type ServerDal interface {
	AddPeer(ctx context.Context, peer *entities.Peer) error
}

type ClientDal interface {
}
