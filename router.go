package main

import (
	"os"

	_ "github.com/PPSKSY-Cluster/backend/docs"
	"github.com/PPSKSY-Cluster/backend/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func InitRouter() (*fiber.App, error) {
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
	api.Get("/docs/*", docsHandler())

	userRoutes := api.Group("/users")
	handlers.InitUserHandlers(userRoutes)

	cResourceRoutes := api.Group("/cresources")
	handlers.InitCResourceHandlers(cResourceRoutes)

	return router, nil
}

// @Description  The route that serves the swagger documentation
// @Tags         general
// @Produce      html
// @Success      200  {html} html
// @Router       /api/docs/ [get]
func docsHandler() func(*fiber.Ctx) error {
	return swagger.HandlerDefault
}

// @Description  Ping route to act as healthcheck
// @Tags         general
// @Produce      json
// @Success      200  {string} string
// @Router       /api/ping [get]
func pingHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}
}
