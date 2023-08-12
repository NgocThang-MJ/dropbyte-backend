package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/liquiddev99/dropbyte-backend/api"
	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config file", err)
	}

	dbpool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer dbpool.Close()

	query := db.New(dbpool)

	server, err := api.NewServer(config, query)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	server.Start(config.ServerAddress)
}
