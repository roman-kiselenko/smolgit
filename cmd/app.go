package cmd

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
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
	exitSig  chan (os.Signal)
}

func New(version string, configPath *string, exitchnl, signchnl chan (os.Signal)) (*App, error) {
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

func (a *App) Run(staticFiles embed.FS) error {
	logger.Info("version", "version", a.Config.Version)
	if !a.Config.ServerDisabled {
		if err := a.initWebServer(staticFiles); err != nil {
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
		a.exitSig <- code
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

func (a *App) initWebServer(staticFiles embed.FS) error {
	logger.Info("initialize web server", "addr", a.Config.ServerAddr)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	r, err := route.New(router, a.Config)
	if err != nil {
		return err
	}
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"title": "Not Found"})
	})
	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/api/auth_disabled", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": a.Config.ServiceAuthDisabled,
		})
	})
	router.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		fileServer := http.FileServer(http.FS(staticFiles))
		c.Request.URL.Path = "/dist" + c.Request.URL.Path
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/api/repos", r.Repos)
	// router.GET("/repo/log/:user/:path", r.Log)
	// router.GET("/repo/files/:user/:path", r.Files)
	// router.GET("/repo/refs/:user/:path", r.Refs)

	addr := a.Config.ServerAddr
	go func() {
		logger.Info("start server", "brand", a.Config.ServerBrand, "address", addr)
		if err := router.Run(addr); err != nil {
			logger.Error("cant run ssh server", "error", err)
		}
	}()
	return nil
}
