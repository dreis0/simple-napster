package dal

import (
	"context"
	"simple-napster/entities"
)

type Dal interface {
	AddPeer(ctx context.Context, peer *entities.Peer) error
}
