package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	mdb, err := InitMongoDB()
	if err != nil {
		return err
	}

	router, err := InitRouter(mdb)
	if err != nil {
		return err
	}

	router.Listen(":3000")

	return nil
}
