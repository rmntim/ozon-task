package inmemory

import (
	"context"
	"github.com/rmntim/ozon-task/internal/models"
)

type Storage struct{}

func (s Storage) CreateUser(ctx context.Context, username string, email string, password string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) CreatePost(ctx context.Context, title string, content string, authorId uint) (*models.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) CreateComment(ctx context.Context, content string, authorId uint, postId uint, parentCommentId *uint) (*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetUsers(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetPostById(ctx context.Context, id uint) (*models.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetPosts(ctx context.Context, limit int, offset int) ([]*models.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetCommentById(ctx context.Context, id uint) (*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetComments(ctx context.Context, limit int, offset int) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) ToggleComments(ctx context.Context, postId uint, userId uint) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetPostsFromUser(ctx context.Context, userId uint) ([]*models.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetReplies(ctx context.Context, commentId uint) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetCommentsForPost(ctx context.Context, postId uint) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func New() *Storage {
	return &Storage{}
}
