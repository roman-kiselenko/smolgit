package main

import (
	"context"
	"log/slog"
	"os"

	"smollgit/cmd"
)

var Version = "dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := cmd.Server(logger, Version)

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("failed to start: %w", "error", err)
		os.Exit(1)
	}
}
