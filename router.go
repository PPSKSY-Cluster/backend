package main

import (
	"os"

	"github.com/PPSKSY-Cluster/backend/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func InitRouter(mdb *mdb) (*fiber.App, error) {
	router := fiber.New()

	// define CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CLIENT_URL"),
		AllowHeaders: "Content-Type",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// define api routes
	api := router.Group("/api")
	api.Get("/ping", pingHandler())
	api.Get("/docs/*", swagger.HandlerDefault)

	var userRoutes fiber.Router = api.Group("/users")
	handlers.InitUserHandlers(userRoutes)

	return router, nil
}

func pingHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}
}
