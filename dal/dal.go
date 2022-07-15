package dal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"simple-napster/entities"
)

type dalImpl struct {
	db *bun.DB
}

var _ ServerDal = (*dalImpl)(nil)
var _ ClientDal = (*dalImpl)(nil)

func NewDal(config *Config) ServerDal {

	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", config.Hostname, config.Port)),
		pgdriver.WithUser(config.Username),
		pgdriver.WithPassword(config.Password),
		pgdriver.WithDatabase(config.Database),
		pgdriver.WithInsecure(true),
	)
	sqldb := sql.OpenDB(conn)

	return &dalImpl{
		db: bun.NewDB(sqldb, pgdialect.New()),
	}
}

func (dal *dalImpl) AddPeer(ctx context.Context, peer *entities.Peer) error {
	_, err := dal.db.NewInsert().Model(peer).Exec(ctx)
	return err
}

func (dal *dalImpl) GetPeers(ctx context.Context) ([]*entities.Peer, error) {
	result := []*entities.Peer{}
	err := dal.db.NewSelect().Model(&result).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dal *dalImpl) UpdatePeer(ctx context.Context, peer *entities.Peer) error {
	_, err := dal.db.NewUpdate().Model(peer).Where("id = ?", peer.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dal *dalImpl) AddFilesIfNew(ctx context.Context, files []*entities.File) error {
	_, err := dal.db.NewInsert().
		Model(files).
		On("CONFLICT (name) DO UPDATE").
		Set("title = EXCLUDED.name").
		Exec(ctx)

	return err
}

func (dal *dalImpl) AddFilesToPeer(ctx context.Context, peer *entities.Peer, files []*entities.File) error {
	peerFiles := make([]*entities.FilePeer, len(files))
	for i, f := range files {
		peerFiles[i] = &entities.FilePeer{PeerId: peer.ID, FileId: f.ID}
	}

	_, err := dal.db.NewInsert().
		Model(peerFiles).
		Ignore().
		Exec(ctx)

	return err
}
