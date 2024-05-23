package main

import (
	"log/slog"
	"os"

	"github.com/rmntim/ozon-task/intenal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg, dbCfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	_ = dbCfg
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
