package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/urfave/cli"

	"git-pebble/cmd"
	"git-pebble/pkg/config"

	"github.com/kelseyhightower/envconfig"
)

var Version = "dev"

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	slog.SetDefault(logger)

	var s config.Config
	if err := envconfig.Process("", &s); err != nil {
		slog.Error("cant parse config", "error", err)
		os.Exit(1)
	}
	app := cli.NewApp()
	app.Name = "gitPebble"
	app.Usage = "a super small and efficient git repository management tool"
	app.Version = fmt.Sprintf("%s-%s", runtime.Version(), Version)
	app.Commands = []cli.Command{
		cmd.Server(logger, s.Database, s.Address),
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to start: %w", "error", err)
		os.Exit(1)
	}
}
