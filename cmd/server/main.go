package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rmntim/ozon-task/graph"
	"github.com/rmntim/ozon-task/graph/resolver"
	"github.com/rmntim/ozon-task/internal/config"
	"github.com/rmntim/ozon-task/internal/lib/logger/sl"
	loggerMw "github.com/rmntim/ozon-task/internal/server/middleware/logger"
	"github.com/rmntim/ozon-task/internal/storage"
	"log/slog"
	"net/http"
	"os"
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

	db, err := storage.New(cfg.Storage, dbCfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	gqlHandler := handler.NewDefaultServer(graph.NewExecutableSchema(resolver.New(db, log)))

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("Ozon Task", "/query"))
	mux.Handle("/query", gqlHandler)

	handlerWithMw := loggerMw.New(log)(mux)
	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      handlerWithMw,
		WriteTimeout: cfg.Server.Timeout,
		ReadTimeout:  cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}

	log.Info("server stopped")
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
