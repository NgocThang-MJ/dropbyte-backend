package gapi

import (
	"log"

	db "github.com/liquiddev99/dropbyte-backend/db/sqlc"
	"github.com/liquiddev99/dropbyte-backend/pb"
	"github.com/liquiddev99/dropbyte-backend/token"
	"github.com/liquiddev99/dropbyte-backend/util"
)

type Server struct {
	pb.UnimplementedDropbyteServer
	config util.Config
	db     *db.Queries
	token  token.Token
}

// Create a new gRPC server
func NewServer(config util.Config, db *db.Queries) (*Server, error) {
	token, err := token.NewMaker(config.SymmetricKey)
	if err != nil {
		log.Fatal("Cannot create token maker")
	}
	server := &Server{config: config, db: db, token: token}

	return server, nil
}
