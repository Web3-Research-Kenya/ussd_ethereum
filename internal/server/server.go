package server

import (
	"github.com/gofiber/fiber/v2"

	"ussd_ethereum/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ussd_ethereum",
			AppName:      "ussd_ethereum",
		}),

		db: database.New(),
	}

	return server
}
