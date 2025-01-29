package handlers

import (
	"ussd_ethereum/internal/database"
	"ussd_ethereum/internal/eth"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Handler struct {
	DB     database.Service
	Tree   *MenuTree
	Dat    Data
	Client *ethclient.Client
}

func NewHandler(db database.Service) *Handler {
	return &Handler{
		Tree:   NewMenuTree(),
		DB:     db,
		Dat:    Data{},
		Client: eth.Connect(),
	}
}
