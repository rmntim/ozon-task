package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rmntim/ozon-task/graph"
	"github.com/rmntim/ozon-task/intenal/config"
	"github.com/rmntim/ozon-task/intenal/lib/logger/sl"
	"github.com/rmntim/ozon-task/intenal/storage"
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

	_ = db

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))

	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", srv)

	if err := http.ListenAndServe(cfg.Server.Address, nil); err != nil {
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
