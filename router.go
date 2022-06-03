package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
)

func InitRouter(mdb *mdb) (*fiber.App, error) {
	router := fiber.New()

	// define CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CLIENT_URL"),
		AllowHeaders: "Content-Type",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// define routes
	router.Get("/api/ping", pingHandler())

	return router, nil
}

func pingHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}
}
