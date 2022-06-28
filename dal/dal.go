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

type DalImpl struct {
	db *bun.DB
}

func NewDal(config *Config) Dal {

	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", config.Hostname, config.Port)),
		pgdriver.WithUser(config.Username),
		pgdriver.WithPassword(config.Password),
		pgdriver.WithDatabase(config.Database),
		pgdriver.WithInsecure(true),
	)
	sqldb := sql.OpenDB(conn)

	return &DalImpl{
		db: bun.NewDB(sqldb, pgdialect.New()),
	}
}

func (dal *DalImpl) AddPeer(ctx context.Context, peer *entities.Peer) error {
	_, err := dal.db.NewInsert().Model(peer).Exec(ctx)
	return err
}
