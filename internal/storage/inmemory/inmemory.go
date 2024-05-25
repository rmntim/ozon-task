package inmemory

import (
	"context"
	"github.com/rmntim/ozon-task/internal/models"
	"github.com/rmntim/ozon-task/internal/server"
	"sync/atomic"
	"time"
)

type user struct {
	id           uint64
	username     string
	email        string
	passwordHash []byte
}

type post struct {
	id                uint64
	title             string
	content           string
	createdAt         time.Time
	authorId          uint64
	commentsAvailable bool
}

type comment struct {
	id              uint64
	content         string
	authorId        uint64
	createdAt       time.Time
	postId          uint64
	parentCommentId *uint
}

type Storage struct {
	users    Map[uint64, *user]
	usersSeq atomic.Uint64

	posts    Map[uint64, *post]
	postsSeq atomic.Uint64

	comments    Map[uint64, *comment]
	commentsSeq atomic.Uint64
}

func New() *Storage {
	return &Storage{
		users:    Map[uint64, *user]{},
		posts:    Map[uint64, *post]{},
		comments: Map[uint64, *comment]{},
	}
}

func (s *Storage) CreateUser(ctx context.Context, username string, email string, password string) (*models.User, error) {
	id := s.usersSeq.Load()

	user := &user{
		id:           id,
		username:     username,
		email:        email,
		passwordHash: []byte(password + "_hashed"), // again, don't care about security here
	}

	s.users.Store(id, user)
	s.usersSeq.Add(1)

	postsIds := make([]uint, 0)
	s.posts.Range(func(id uint64, p *post) bool {
		if p.authorId == id {
			postsIds = append(postsIds, uint(p.id))
		}
		return true
	})

	return &models.User{
		ID:       uint(user.id),
		Username: user.username,
		Email:    user.email,
		PostsIDs: postsIds,
	}, nil
}

func (s *Storage) CreatePost(ctx context.Context, title string, content string, authorId uint) (*models.Post, error) {
	id := s.postsSeq.Load()

	post := &post{
		id:                id,
		title:             title,
		content:           content,
		createdAt:         time.Now(),
		authorId:          uint64(authorId),
		commentsAvailable: true,
	}

	s.posts.Store(id, post)
	s.postsSeq.Add(1)

	commentsIds := make([]uint, 0)
	s.comments.Range(func(id uint64, c *comment) bool {
		if c.postId == id {
			commentsIds = append(commentsIds, uint(c.id))
		}
		return true
	})

	return &models.Post{
		ID:          uint(post.id),
		Title:       post.title,
		CreatedAt:   post.createdAt,
		Content:     post.content,
		AuthorID:    uint(post.authorId),
		CommentsIDs: commentsIds,
	}, nil
}

func (s *Storage) CreateComment(ctx context.Context, content string, authorId uint, postId uint, parentCommentId *uint) (*models.Comment, error) {
	id := s.commentsSeq.Load()

	newComment := &comment{
		id:              id,
		content:         content,
		authorId:        uint64(authorId),
		createdAt:       time.Now(),
		postId:          uint64(postId),
		parentCommentId: parentCommentId,
	}

	s.comments.Store(id, newComment)
	s.commentsSeq.Add(1)

	commentsIds := make([]uint, 0)
	s.comments.Range(func(id uint64, c *comment) bool {
		if c.postId == id {
			commentsIds = append(commentsIds, uint(c.id))
		}
		return true
	})

	return &models.Comment{
		ID:              uint(newComment.id),
		Content:         newComment.content,
		AuthorID:        uint(newComment.authorId),
		CreatedAt:       newComment.createdAt,
		PostID:          uint(newComment.postId),
		ParentCommentID: parentCommentId,
		RepliesIDs:      commentsIds,
	}, nil
}

func (s *Storage) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	user, ok := s.users.Load(uint64(id))
	if !ok {
		return nil, server.ErrUserNotFound
	}

	postsIds := make([]uint, 0)
	s.posts.Range(func(id uint64, p *post) bool {
		if p.authorId == id {
			postsIds = append(postsIds, uint(p.id))
		}
		return true
	})

	return &models.User{
		ID:       uint(user.id),
		Username: user.username,
		Email:    user.email,
		PostsIDs: postsIds,
	}, nil
}

func (s *Storage) GetUsers(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	users := make([]*models.User, limit)

	for i := 0; i < limit; i++ {
		user, ok := s.users.Load(uint64(offset + i))
		if !ok {
			break
		}

		postsIds := make([]uint, 0)
		s.posts.Range(func(id uint64, p *post) bool {
			if p.authorId == id {
				postsIds = append(postsIds, uint(p.id))
			}
			return true
		})

		users[i] = &models.User{
			ID:       uint(user.id),
			Username: user.username,
			Email:    user.email,
			PostsIDs: postsIds,
		}
	}

	return users, nil
}

func (s *Storage) GetPostById(ctx context.Context, id uint) (*models.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetPosts(ctx context.Context, limit int, offset int) ([]*models.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetCommentById(ctx context.Context, id uint) (*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetComments(ctx context.Context, limit int, offset int) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) ToggleComments(ctx context.Context, postId uint, userId uint) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetPostsFromUser(ctx context.Context, userId uint) ([]*models.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetReplies(ctx context.Context, commentId uint) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetCommentsForPost(ctx context.Context, postId uint) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}
