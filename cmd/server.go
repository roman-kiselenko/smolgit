package cmd

import (
	"log/slog"

	sloggin "github.com/samber/slog-gin"
	altsrc "github.com/urfave/cli-altsrc/v3"

	"smollgit/internal/route"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"
)

func Server(log *slog.Logger, version string) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start server",
		Description: "smollgit server handles all stuff for you",
		Version:     version,
		Action:      runServer(log),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Sources: altsrc.YAML("server.address", []string{"./config.yaml"}...),
			},
			&cli.StringMapFlag{},
		},
	}
}

func runServer(log *slog.Logger) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		gin.SetMode(gin.ReleaseMode)
		router := gin.New()
		router.Use(sloggin.New(log))
		router.Use(gin.Recovery())

		r, err := route.New()
		if err != nil {
			return err
		}

		router.GET("/", r.Index)

		addr := ctx.String("address")
		log.Info("Start smollgit instance...", "address", addr)
		return router.Run(addr)
	}
}
