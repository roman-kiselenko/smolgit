package cmd

import (
	"log/slog"
	"os"

	"git-pebble/internal/db"
	"git-pebble/internal/route"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/urfave/cli"
)

func Server(log *slog.Logger, dbPath, address string) cli.Command {
	return cli.Command{
		Name:        "server",
		Usage:       "Start server",
		Description: "gitPebble server handles all stuff for you",
		Action:      runServer(log, dbPath, address),
	}
}

func runServer(log *slog.Logger, dbPath, address string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.New(log))
	router.Use(gin.Recovery())

	database, err := db.New(dbPath)
	if err != nil {
		slog.Error("cant connect sqlite", "error", err)
		os.Exit(1)
	}
	r, err := route.New(database)
	if err != nil {
		slog.Error("cant init routes", "error", err)
		os.Exit(1)
	}

	router.GET("/", r.Index)

	log.Info("Start gitPebble instance...", "address", address)
	return router.Run(address)
}
