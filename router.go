package main

import "github.com/gofiber/fiber/v2"

func InitRouter() (*fiber.App, error) {
	router := fiber.New()

	router.Get("/ping", pingHandler())

	return router, nil
}

func pingHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}
}
