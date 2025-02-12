package server

import (
	"fmt"
	"ussd_ethereum/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {

	h := handlers.NewHandler(s.db)
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Post("/callback", h.CallbackHandler)
	s.App.Post("/reports", Reports)

	s.App.Get("/health", s.healthHandler)

}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func Reports(c *fiber.Ctx) error {
	fmt.Printf("REPORT: %+v\n", string(c.Body()))
	return nil
}
