package main

import (
	"log"

	"smolgit/cmd"
)

var Version = "dev"

func main() {
	app, err := cmd.New(Version)
	if err != nil {
		log.Fatalf("failed to init app: %s", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("failed to start app: %s", err)
	}
}
