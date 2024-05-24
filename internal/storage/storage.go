package storage

import (
	"context"
	"fmt"
	"github.com/rmntim/ozon-task/graph/model"
	"github.com/rmntim/ozon-task/internal/config"
	"github.com/rmntim/ozon-task/internal/storage/postgres"
)

type Storage interface {
	CreateUser(ctx context.Context, username string, email string, password string) (*model.User, error)
	CreatePost(ctx context.Context, title string, content string, authorId uint) (*model.Post, error)
	CreateComment(ctx context.Context, content string, authorId uint, postId uint, parentCommentId *uint) (*model.Comment, error)
	GetUserById(ctx context.Context, id uint) (*model.User, error)
	GetUsers(ctx context.Context, limit int, offset int) ([]*model.User, error)
	GetPostById(ctx context.Context, id uint) (*model.Post, error)
	GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, error)
	GetCommentById(ctx context.Context, id uint) (*model.Comment, error)
	GetComments(ctx context.Context, limit int, offset int) ([]*model.Comment, error)
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
