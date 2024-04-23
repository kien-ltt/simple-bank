package gapi

import (
	"fmt"

	db "github.com/kien-ltt/simplebank/db/sqlc"
	"github.com/kien-ltt/simplebank/pb"
	"github.com/kien-ltt/simplebank/token"
	"github.com/kien-ltt/simplebank/util"
	"github.com/kien-ltt/simplebank/worker"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		tokenMaker:      tokenMaker,
		store:           store,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
