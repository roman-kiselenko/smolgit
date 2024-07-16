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
			&cli.StringFlag{
				Name:    "ssh_addr",
				Usage:   "Address for SSH server",
				Value:   ":3081",
				Sources: altsrc.YAML("ssh.addr", configs...),
			},
			&cli.BoolFlag{
				Name:    "log_color",
				Usage:   "Color output of logs, works only with text logs",
				Value:   true,
				Sources: altsrc.YAML("log.color", configs...),
			},
			&cli.BoolFlag{
				Name:    "log_json",
				Usage:   "Output logs in json format",
				Value:   false,
				Sources: altsrc.YAML("log.json", configs...),
			},
			&cli.StringFlag{
				Name:    "log_level",
				Usage:   "Set log level (INFO, DEBUG, WARN, TRACE)",
				Value:   "DEBUG",
				Sources: altsrc.YAML("log.level", configs...),
			},
			&cli.StringFlag{
				Name:    "server_addr",
				Usage:   "Address for Web server",
				Value:   ":3080",
				Sources: altsrc.YAML("server.addr", configs...),
			},
			&cli.StringFlag{
				Name:    "server_brand",
				Usage:   "Header used in web",
				Value:   "smolgit",
				Sources: altsrc.YAML("server.brand", configs...),
			},
			&cli.StringFlag{
				Name:    "git_base",
				Usage:   "Format 'git clone ...' link",
				Value:   "git@my-git-server.lan",
				Sources: altsrc.YAML("git.base", configs...),
			},
			&cli.StringFlag{
				Name:    "git_path",
				Usage:   "Set directory for git repositories",
				Value:   "./tmp/smolgit",
				Sources: altsrc.YAML("git.path", configs...),
			},
			&cli.StringFlag{
				Name:    "db_path",
				Usage:   "Path to save database file",
				Value:   "./database.sqlite",
				Sources: altsrc.YAML("db.path", configs...),
			},
			&cli.StringFlag{
				Name:    "root_login",
				Usage:   "Admin user will be create at first start",
				Sources: altsrc.YAML("root.login", configs...),
			},
			&cli.StringFlag{
				Name:    "root_password",
				Usage:   "Admin password will be create at first start",
				Sources: altsrc.YAML("root.password", configs...),
			},
			&cli.StringSliceFlag{
				Name:    "root_keys",
				Usage:   "Admin ssh-keys will be create at first start",
				Sources: altsrc.YAML("root.keys", configs...),
			},
		},
	}
}

// TODO move all check to PreStart
func initApp(ctx *cli.Context) error {
	logger := initLogger(ctx)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())
	if fileName, ok := checkConfigFiles(); ok {
		logger.Info("found config file", "config", fileName)
	}

	logger.Info("version", "version", ctx.Command.Version)
	logger.Info("initialize sqlite database", "db_path", ctx.String("db_path"))
	database, err := db.New(logger, ctx)
	if err != nil {
		logger.Error("cant connect sqlite", "error", err)
		return err
	}
	gitPath := ctx.String("git_path")

	r, err := route.New(logger, router, database, gitPath, ctx.String("git_base"))
	if err != nil {
		logger.Error("cant create routes", "error", err)
		return err
	}

	router.GET("/", r.Index)
	router.GET("/css/pico.min.css", r.ExternalStyle)
	router.GET("/css/style.css", r.Style)
	router.GET("/repo/log/:user/:path", r.Log)
	router.GET("/repo/files/:user/:path", r.Files)
	router.GET("/repo/refs/:user/:path", r.Refs)
	router.GET("/users", r.Users)
	router.POST("/user", r.PostUser)
	router.POST("/repo", r.PostRepo)
	router.GET("/create/user", r.CreateUser)
	router.GET("/create/repo", r.CreateRepo)

	addr := ctx.String("server_addr")

	logger.Info("initialize ssh server", "addr", ctx.String("ssh_addr"))
	sshServer, err := ssh.New(logger, database, ctx)
	if err != nil {
		logger.Error("cant run ssh server", "error", err)
		return err
	}

	go func() {
		logger.Info("starting SSH server", "addr", ctx.String("ssh_addr"))
		if err := sshServer.ListenAndServe(); err != nil {
			logger.Error("cant run ssh server", "error", err)
		}
	}()

	logger.Info("start server", "brand", ctx.String("server_brand"), "address", addr)
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
	logger.Info("set loglevel", "level", level)

	return logger
}
