package resolver

import (
	"github.com/rmntim/ozon-task/graph"
	"github.com/rmntim/ozon-task/internal/storage"
	"log/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db  storage.Storage
	log *slog.Logger
}

func New(db storage.Storage, log *slog.Logger) graph.ResolverRoot {
	return &Resolver{
		db:  db,
		log: log,
	}
}
