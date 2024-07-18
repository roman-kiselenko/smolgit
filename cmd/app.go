package cmd

import (
	"fmt"
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
)

var (
	k         = koanf.New(".")
	logOutput = os.Stdout
	logger    *slog.Logger
	cfg       config.Config
)

type App struct {
	Config   *config.Config
	signchnl chan (os.Signal)
	exitSig  chan (int)
}

func New(version string, configPath *string, exitchnl chan (int), signchnl chan (os.Signal)) (*App, error) {
	app := &App{exitSig: exitchnl, signchnl: signchnl}
	f := file.Provider(*configPath)
	cfg.Version = version
	k = koanf.New(".")
	if err := k.Load(f, yaml.Parser()); err != nil {
		return app, err
	}
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true}); err != nil {
		return app, err
	}
	app.Config = &cfg
	logger = initLogger()
	return app, nil
}

func (a *App) Run() error {
	logger.Info("version", "version", a.Config.Version)

	if !a.Config.ServerDisabled {
		if err := a.initWebServer(); err != nil {
			return fmt.Errorf("cant run web server %w", err)
		}
	}

	logger.Info("initialize ssh server", "addr", a.Config.SSHAddr)
	sshServer, err := ssh.New(a.Config)
	if err != nil {
		return fmt.Errorf("cant run ssh server %w", err)
	}

	go func() {
		logger.Info("starting SSH server", "addr", a.Config.SSHAddr)
		if err := sshServer.ListenAndServe(); err != nil {
			logger.Error("cant run ssh server", "error", err)
		}
	}()

	go func() {
		code := <-a.signchnl
		logger.Info("os signal received", "signal", code)
		if err := sshServer.Close(); err != nil {
			logger.Error("cant stop ssh server", "error", err)
		}
		a.exitSig <- 0
	}()

	return nil
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

func (a *App) initWebServer() error {
	logger.Info("initialize web server", "addr", a.Config.ServerAddr)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	if cfg.ServerAuthEnabled {
		logger.Info("web auth", "enabled", a.Config.ServerAuthEnabled)
		accounts := gin.Accounts{}
		for _, a := range cfg.ServerAuthAccounts {
			accounts[a["login"]] = a["password"]
		}
		router.Use(gin.BasicAuth(accounts))
	}

	r, err := route.New(router, a.Config)
	if err != nil {
		return err
	}

	router.GET("/", r.Index)
	router.GET("/css/pico.min.css", r.ExternalStyle)
	router.GET("/css/style.css", r.Style)
	router.GET("/repo/log/:user/:path", r.Log)
	router.GET("/repo/files/:user/:path", r.Files)
	router.GET("/repo/refs/:user/:path", r.Refs)

	addr := a.Config.ServerAddr
	go func() {
		logger.Info("start server", "brand", a.Config.ServerBrand, "address", addr)
		if err := router.Run(addr); err != nil {
			logger.Error("cant run ssh server", "error", err)
		}
	}()
	return nil
}
