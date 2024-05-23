package storage

import (
	"fmt"
	"github.com/rmntim/ozon-task/internal/config"
	"github.com/rmntim/ozon-task/internal/storage/postgres"
)

type Storage interface {
}

// New creates new storage instance, depending on storage type.
func New(storageType string, dbCfg *config.DBConfig) (Storage, error) {
	switch storageType {
	case "postgres":
		return postgres.New(dbCfg.Username, dbCfg.Password, dbCfg.Address, dbCfg.Database)
	case "memory":
		storage := make(map[string]string, 10)
		return storage, nil
	}

	return nil, fmt.Errorf("unknown storage type: %s", storageType)
}
