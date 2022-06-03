package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load("./.env")
	if err != nil {
		return err
	}

	mdb, err := InitMongoDB()
	if err != nil {
		return err
	}

	router, err := InitRouter(mdb)
	if err != nil {
		return err
	}

	port := os.Getenv("PORT")
	router.Listen(":" + port)

	return nil
}
