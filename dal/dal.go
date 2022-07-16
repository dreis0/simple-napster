package dal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (dal *dalImpl) AddFileToPeeerWithFilename(ctx context.Context, peer *entities.Peer, filename string) error {
	file := &entities.File{}
	err := dal.db.NewSelect().Model(file).Where("name = ?", filename).Scan(ctx)
	if err != nil {
		return err
	}

	return dal.AddFilesToPeer(ctx, peer, []*entities.File{file})
}

func (dal *dalImpl) GetPeerById(ctx context.Context, id uuid.UUID) (*entities.Peer, error) {
	peer := &entities.Peer{}
	err := dal.db.NewSelect().Model(peer).Where("id = ?", id).Scan(ctx)

	return peer, err
}

func (dal *dalImpl) DeletePeerAndFiles(ctx context.Context, id uuid.UUID) error {
	_, err := dal.db.NewDelete().Model(&entities.Peer{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = dal.db.NewDelete().Model(&entities.FilePeer{}).Where("peer_id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (dal *dalImpl) GetAllPeersWithFile(ctx context.Context, filename string) ([]*entities.Peer, error) {
	exists, err := dal.db.NewSelect().Model(&entities.File{}).Where("name = ?", filename).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("file not found")
	}

	peers := []*entities.Peer{}
	err = dal.db.NewSelect().Model(peers).
		Join("JOIN file_peer fp on fp.peer_id = p.id").
		Join("JOIN files f on f.id = fp.id").
		Where("active = 1").
		Where("f.name = ?", filename).
		Scan(ctx)

	return peers, nil
}
