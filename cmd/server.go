package cmd

import (
	"log/slog"
	"os"
	"time"

	"smolgit/pkg/config"
	"smolgit/pkg/route"
	"smolgit/pkg/ssh"

	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/lmittmann/tint"
	sloggin "github.com/samber/slog-gin"
)

var (
	k         = koanf.New(".")
	logOutput = os.Stdout
	logger    *slog.Logger
	cfg       config.Config
)

type App struct {
	Config *config.Config
}

func New(version string) (*App, error) {
	app := &App{}
	f := file.Provider("./config.yaml")
	cfg.Version = version
	k = koanf.New(".")
	if err := k.Load(f, yaml.Parser()); err != nil {
		return app, err
	}
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true}); err != nil {
		return app, err
	}
	app.Config = &cfg
	err := f.Watch(func(_ interface{}, _ error) {
		k = koanf.New(".")
		k.Load(f, yaml.Parser())
		k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true})
		logger.Info("reloading config...", "cfg", "./config.yaml")
		app.Config = &cfg
	})
	return app, err
}

// TODO move all check to PreStart
func (a *App) Sleep() error {
	logger := initLogger()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())
	logger.Info("version", "version", a.Config.Version)

	r, err := route.New(logger, router, a.Config)
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
	router.POST("/repo", r.PostRepo)
	router.GET("/create/repo", r.CreateRepo)

	addr := a.Config.ServerAddr

	logger.Info("initialize ssh server", "addr", addr)
	sshServer, err := ssh.New(logger, a.Config)
	if err != nil {
		logger.Error("cant run ssh server", "error", err)
		return err
	}

	go func() {
		logger.Info("starting SSH server", "addr", a.Config.SSHAddr)
		if err := sshServer.ListenAndServe(); err != nil {
			logger.Error("cant run ssh server", "error", err)
		}
	}()

	logger.Info("start server", "brand", a.Config.ServerBrand, "address", addr)
	return router.Run(addr)
}

func initLogger() *slog.Logger {
	level := new(slog.LevelVar)
	handler := &slog.HandlerOptions{
		Level: level,
	}
	if cfg.LogJSON {
		logger = slog.New(slog.NewJSONHandler(logOutput, handler))
	} else {
		if cfg.LogColor {
			logger = slog.New(tint.NewHandler(logOutput, &tint.Options{
				Level:      level,
				TimeFormat: time.Kitchen,
			}))
		} else {
			logger = slog.New(slog.NewTextHandler(logOutput, handler))
		}
	}

	slog.SetDefault(logger)
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		level.Set(slog.LevelDebug)
	}
	logger.Info("set loglevel", "level", level)

	return logger
}
