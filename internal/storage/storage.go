package storage

import (
	"context"
	"fmt"
	"github.com/rmntim/ozon-task/internal/config"
	"github.com/rmntim/ozon-task/internal/models"
	"github.com/rmntim/ozon-task/internal/storage/inmemory"
	"github.com/rmntim/ozon-task/internal/storage/postgres"
)

type Storage interface {
	CreateUser(ctx context.Context, username string, email string, password string) (*models.User, error)
	CreatePost(ctx context.Context, title string, content string, authorId uint) (*models.Post, error)
	CreateComment(ctx context.Context, content string, authorId uint, postId uint, parentCommentId *uint) (*models.Comment, error)
	GetUserById(ctx context.Context, id uint) (*models.User, error)
	GetUsers(ctx context.Context, limit int, offset int) ([]*models.User, error)
	GetPostById(ctx context.Context, id uint) (*models.Post, error)
	GetPosts(ctx context.Context, limit int, offset int) ([]*models.Post, error)
	GetCommentById(ctx context.Context, id uint) (*models.Comment, error)
	GetComments(ctx context.Context, limit int, offset int) ([]*models.Comment, error)
	ToggleComments(ctx context.Context, postId uint, userId uint) (bool, error)
	GetPostsFromUser(ctx context.Context, userId uint) ([]*models.Post, error)
	GetReplies(ctx context.Context, commentId uint) ([]*models.Comment, error)
	GetCommentsForPost(ctx context.Context, postId uint) ([]*models.Comment, error)
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
		return inmemory.New(), nil
	}
	return nil, fmt.Errorf("unknown storage type: %s", storageType)
}
