package main

import (
	"fmt"
	"os"

	"github.com/PPSKSY-Cluster/backend/api"
	"github.com/PPSKSY-Cluster/backend/auth"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/joho/godotenv"
)

// @title PPSKSY-Cluster API
// @version 1.0
// @description This is the API for the PPSKSY-Cluster Webapplication
// @license.name MIT
// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load("./.env"); err != nil {
		return err
	}

	if err := db.InitDB(); err != nil {
		return err
	}

	if err := auth.InitAuth(); err != nil {
		return err
	}

	router, err := api.InitRouter()
	if err != nil {
		return err
	}

	port := os.Getenv("PORT")
	router.Listen(":" + port)

	return nil
}
