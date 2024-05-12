package server

import (
	"github.com/gofiber/fiber/v2"

	"health-record/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "health-record",
			AppName:      "health-record",
		}),

		db: database.New(),
	}

	return server
}
