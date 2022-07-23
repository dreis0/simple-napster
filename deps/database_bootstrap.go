package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"os"
	"simple-napster/dal"
	"simple-napster/entities"
	"simple-napster/utils"
)

func main() {
	file, err := utils.GetArgument(os.Args, "env")
	if err != nil {
		panic(err)
	}
	err = godotenv.Load(file)
	if err != nil {
		panic(err)
	}

	createDatabase()
	config := dal.FromEnv()

	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(config.Url),
		pgdriver.WithUser(config.Username),
		pgdriver.WithPassword(config.Password),
		pgdriver.WithDatabase(config.Database),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()

	//enableUUID(db)
	createPeersTable(db, ctx)
	createFilesTable(db, ctx)
	createFilePeerTable(db, ctx)
}

func createDatabase() {
	config := dal.FromEnv()
	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf(config.Url)),
		pgdriver.WithUser(config.Username),
		pgdriver.WithPassword(config.Password),
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
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err != nil {
		panic(err)
	}
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

func createFilesTable(db *bun.DB, ctx context.Context) {
	_, err := db.NewCreateTable().
		Model(&entities.File{}).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		panic(err)
	}
}

func createFilePeerTable(db *bun.DB, ctx context.Context) {
	_, err := db.NewCreateTable().
		Model(&entities.FilePeer{}).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		panic(err)
	}
}
