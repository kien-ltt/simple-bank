package gapi

import (
	"fmt"

	db "github.com/kien-ltt/simplebank/db/sqlc"
	"github.com/kien-ltt/simplebank/pb"
	"github.com/kien-ltt/simplebank/token"
	"github.com/kien-ltt/simplebank/util"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		store:      store,
	}

	return server, nil
}
