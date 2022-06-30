package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"simple-napster/entities"
)

func main() {
	createDatabase()

	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("localhost:5432")),
		pgdriver.WithUser("postgres"),
		pgdriver.WithPassword("postgres"),
		pgdriver.WithDatabase("napster"),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()

	enableUUID(db)
	createPeersTable(db, ctx)
}

func createDatabase() {
	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("localhost:5432")),
		pgdriver.WithUser("postgres"),
		pgdriver.WithPassword("postgres"),
		pgdriver.WithDatabase("postgres"),
		pgdriver.WithInsecure(true),
	)
	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())

	_, err := db.Exec("DROP DATABASE IF EXISTS napster")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE napster")
	if err != nil {
		panic(err)
	}
}

func enableUUID(db *bun.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
}

func createPeersTable(db *bun.DB, ctx context.Context) {
	_, err := db.NewCreateTable().
		Model(&entities.Peer{}).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		panic(err)
	}
}
