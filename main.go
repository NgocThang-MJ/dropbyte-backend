package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/liquiddev99/dropbyte-backend/api"
	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/gapi"
	"github.com/liquiddev99/dropbyte-backend/pb"
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

	runGinServer(config, query)
	// go runGatewayServer(config, query)
	// runGrpcServer(config, query)
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

// Run gRPC server
func runGrpcServer(config util.Config, query *db.Queries) {
	server, err := gapi.NewServer(config, query)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterDropbyteServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start gRPC server", err)
	}
}

func runGatewayServer(config util.Config, query *db.Queries) {
	server, err := gapi.NewServer(config, query)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterDropbyteHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Cannot start HTTP Gateway server", err)
	}
}

// Run HTTP server
func runGinServer(config util.Config, query *db.Queries) {
	server, err := api.NewServer(config, query)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	log.Println("Starting server at 0.0.0.0:8080")
	server.Start(config.HTTPServerAddress)
}
