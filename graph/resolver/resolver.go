package resolver

import (
	"github.com/rmntim/ozon-task/internal/storage"
	"log/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// TODO: use dataloaders

type Resolver struct {
	db  storage.Storage
	log *slog.Logger
}

// New initializes new resolver with storage instance.
func New(db storage.Storage, log *slog.Logger) *Resolver {
	return &Resolver{db: db, log: log}
}
