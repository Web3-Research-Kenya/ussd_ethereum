package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/logger"
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

	server.db.CreateTable()

	server.Use(logger.New())
	return server
}
