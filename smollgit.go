package main

import (
	"context"
	"log"
	"os"

	"smolgit/cmd"
)

var Version = "dev"

func main() {
	app := cmd.Server(Version)

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("failed to start: %s", err)
	}
}
