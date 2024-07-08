package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	sloggin "github.com/samber/slog-gin"
	altsrc "github.com/urfave/cli-altsrc/v3"

	"smolgit/internal/db"
	"smolgit/internal/route"
	"smolgit/internal/ssh"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"
)

var (
	configs   = []string{"./config.yaml", "./config.yml"}
	logOutput = os.Stdout
)

func Server(version string) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start server",
		Description: "smolgit server handles all stuff for you",
		Version:     version,
		Action:      initApp,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "log_color",
				Value:   true,
				Sources: altsrc.YAML("log.color", configs...),
			},
			&cli.BoolFlag{
				Name:    "log_json",
				Value:   false,
				Sources: altsrc.YAML("log.json", configs...),
			},
			&cli.StringFlag{
				Name:    "log_level",
				Value:   "DEBUG",
				Sources: altsrc.YAML("log.level", configs...),
			},
			&cli.StringFlag{
				Name:    "server_address",
				Value:   ":8081",
				Sources: altsrc.YAML("server.address", configs...),
			},
			&cli.StringFlag{
				Name:    "server_brand",
				Value:   "smolgit",
				Sources: altsrc.YAML("server.brand", configs...),
			},
			&cli.StringFlag{
				Name:    "git_path",
				Value:   "/tmp/smolgit",
				Sources: altsrc.YAML("git.path", configs...),
			},
			&cli.StringFlag{
				Name:    "db_path",
				Value:   "./database.sqlite",
				Sources: altsrc.YAML("db.path", configs...),
			},
		},
	}
}

func initApp(ctx *cli.Context) error {
	logger := initLogger(ctx)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())

	dbPath := ctx.String("db_path")

	database, err := db.New(logger, dbPath)
	if err != nil {
		logger.Error("cant connect sqlite", "error", err)
		return err
	}
	r, err := route.New(logger, router, database)
	if err != nil {
		return err
	}

	router.GET("/", r.Index)
	router.GET("/css/terminal.min.css", r.ExternalStyle)
	router.GET("/css/style.css", r.Style)
	router.GET("/repos", r.Repos)

	addr := ctx.String("server_address")

	if fileName, ok := checkConfigFiles(); ok {
		logger.Info("found config file", "config", fileName)
	}
	logger.Info("init sqlite database", "db_path", dbPath)
	logger.Info("init git directory", "directory", ctx.String("git_path"))
	logger.Info("start server", "brand", ctx.String("server_brand"), "address", addr)

	ssh.Listen(logger, database, ctx)
	return router.Run(addr)
}

func checkConfigFiles() (string, bool) {
	for _, f := range configs {
		if _, err := os.Stat(f); err == nil {
			return f, true
		}
	}
	return "", false
}

func initLogger(ctx *cli.Context) *slog.Logger {
	level := new(slog.LevelVar)
	handler := &slog.HandlerOptions{
		Level: level,
	}
	var logger *slog.Logger
	if ctx.Bool("log_json") {
		logger = slog.New(slog.NewJSONHandler(logOutput, handler))
	} else {
		if ctx.Bool("log_color") {
			logger = slog.New(tint.NewHandler(logOutput, &tint.Options{
				Level:      level,
				TimeFormat: time.Kitchen,
			}))
		} else {
			logger = slog.New(slog.NewTextHandler(logOutput, handler))
		}
	}

	slog.SetDefault(logger)
	if err := level.UnmarshalText([]byte(ctx.String("log_level"))); err != nil {
		level.Set(slog.LevelDebug)
	}

	return logger
}
