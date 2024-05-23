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
		db, err := postgres.New(dbCfg.Username, dbCfg.Password, dbCfg.Address, dbCfg.Database)
		if err != nil {
			return nil, err
		}
		if err := db.Migrate(); err != nil {
			return nil, err
		}
		return db, nil
	case "memory":
		// TODO: improve memory storage
		panic("unimplemented")
	}
	return nil, fmt.Errorf("unknown storage type: %s", storageType)
}
