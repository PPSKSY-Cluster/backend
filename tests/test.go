package tests

import (
	"github.com/PPSKSY-Cluster/backend/api"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func SetupTestApplication() (*fiber.App, error) {
	// Server setup : load env file, setup db, setup router
	if err := godotenv.Load("./.env"); err != nil {
		return nil, err
	}
	if err := db.InitDB(); err != nil {
		return nil, err
	}

	app, err := api.InitRouter()
	if err != nil {
		return nil, err
	}

	return app, nil
}
