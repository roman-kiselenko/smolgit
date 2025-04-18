package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"smolgit/cmd"
	"smolgit/pkg/config"
)

var (
	version    = "dev"
	configPath = flag.String("config", "./config.yaml", "path to config")
)

func main() {
	flag.Parse()

	command := os.Args[1:]
	if len(command) > 0 && command[0] == "config" {
		cfg, err := config.GenerateConfig()
		if err != nil {
			log.Fatalf("failed to generate config: %s", err)
		}
		fmt.Print(string(cfg))
		os.Exit(0)
	}
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	exitchnl := make(chan os.Signal)
	app, err := cmd.New(version, configPath, exitchnl, sigchnl)
	if err != nil {
		log.Fatalf("failed to init app: %s", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("failed to start app: %s", err)
	}
	<-exitchnl
}
