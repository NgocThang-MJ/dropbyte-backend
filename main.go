package main

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
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

	runDbMigration(config.MigrationUrl, config.DatabaseUrl)

	server, err := api.NewServer(config, query)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	log.Println("Starting server at 0.0.0.0:8080")
	server.Start(config.ServerAddress)
}

func runDbMigration(migrationUrl string, dbURL string) {
	migration, err := migrate.New(migrationUrl, dbURL)
	if err != nil {
		log.Fatal("Cannot create new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migration", err)
	}

	log.Println("DB migrated successfully")
}
