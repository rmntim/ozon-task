package inmemory

import (
	"context"
	"github.com/rmntim/ozon-task/internal/models"
	"github.com/rmntim/ozon-task/internal/server"
	"sync/atomic"
	"time"
)

type User struct {
	id           uint64
	username     string
	email        string
	passwordHash []byte
}

type Post struct {
	id                uint64
	title             string
	content           string
	createdAt         time.Time
	authorId          uint64
	commentsAvailable bool
}

type Comment struct {
	id              uint64
	content         string
	authorId        uint64
	createdAt       time.Time
	postId          uint64
	parentCommentId *uint
}

type Storage struct {
	users    Map[uint64, *User]
	usersSeq atomic.Uint64

	posts    Map[uint64, *Post]
	postsSeq atomic.Uint64

	comments    Map[uint64, *Comment]
	commentsSeq atomic.Uint64
}

func New() *Storage {
	return &Storage{
		users:    Map[uint64, *User]{},
		posts:    Map[uint64, *Post]{},
		comments: Map[uint64, *Comment]{},
	}
}

func (s *Storage) CreateUser(ctx context.Context, username string, email string, password string) (*models.User, error) {
	id := s.usersSeq.Load()

	user := &User{
		id:           id,
		username:     username,
		email:        email,
		passwordHash: []byte(password + "_hashed"), // again, don't care about security here
	}

	s.users.Store(id, user)
	s.usersSeq.Add(1)

	postsIds := make([]uint, 0)
	s.posts.Range(func(id uint64, p *Post) bool {
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

	post := &Post{
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
	s.comments.Range(func(id uint64, c *Comment) bool {
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

	comment := &Comment{
		id:              id,
		content:         content,
		authorId:        uint64(authorId),
		createdAt:       time.Now(),
		postId:          uint64(postId),
		parentCommentId: parentCommentId,
	}

	s.comments.Store(id, comment)
	s.commentsSeq.Add(1)

	commentsIds := make([]uint, 0)
	s.comments.Range(func(id uint64, c *Comment) bool {
		if c.postId == id {
			commentsIds = append(commentsIds, uint(c.id))
		}
		return true
	})

	return &models.Comment{
		ID:              uint(comment.id),
		Content:         comment.content,
		AuthorID:        uint(comment.authorId),
		CreatedAt:       comment.createdAt,
		PostID:          uint(comment.postId),
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
	s.posts.Range(func(id uint64, p *Post) bool {
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
		s.posts.Range(func(id uint64, p *Post) bool {
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
	post, ok := s.posts.Load(uint64(id))
	if !ok {
		return nil, server.ErrPostNotFound
	}

	commentsIds := make([]uint, 0)
	s.comments.Range(func(id uint64, c *Comment) bool {
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

func (s *Storage) GetPosts(ctx context.Context, limit int, offset int) ([]*models.Post, error) {
	posts := make([]*models.Post, limit)

	for i := 0; i < limit; i++ {
		post, ok := s.posts.Load(uint64(offset + i))
		if !ok {
			break
		}

		commentsIds := make([]uint, 0)
		s.comments.Range(func(id uint64, c *Comment) bool {
			if c.postId == id {
				commentsIds = append(commentsIds, uint(c.id))
			}
			return true
		})

		posts[i] = &models.Post{
			ID:          uint(post.id),
			Title:       post.title,
			CreatedAt:   post.createdAt,
			Content:     post.content,
			AuthorID:    uint(post.authorId),
			CommentsIDs: commentsIds,
		}
	}

	return posts, nil
}

func (s *Storage) GetCommentById(ctx context.Context, id uint) (*models.Comment, error) {
	comment, ok := s.comments.Load(uint64(id))
	if !ok {
		return nil, server.ErrCommentNotFound
	}

	commentsIds := make([]uint, 0)
	s.comments.Range(func(id uint64, c *Comment) bool {
		if c.postId == id {
			commentsIds = append(commentsIds, uint(c.id))
		}
		return true
	})

	return &models.Comment{
		ID:              uint(comment.id),
		Content:         comment.content,
		AuthorID:        uint(comment.authorId),
		CreatedAt:       comment.createdAt,
		PostID:          uint(comment.postId),
		ParentCommentID: comment.parentCommentId,
		RepliesIDs:      commentsIds,
	}, nil
}

func (s *Storage) GetComments(ctx context.Context, limit int, offset int) ([]*models.Comment, error) {
	comments := make([]*models.Comment, limit)

	for i := 0; i < limit; i++ {
		comment, ok := s.comments.Load(uint64(offset + i))
		if !ok {
			break
		}

		commentsIds := make([]uint, 0)
		s.comments.Range(func(id uint64, c *Comment) bool {
			if c.postId == id {
				commentsIds = append(commentsIds, uint(c.id))
			}
			return true
		})

		comments[i] = &models.Comment{
			ID:              uint(comment.id),
			Content:         comment.content,
			AuthorID:        uint(comment.authorId),
			CreatedAt:       comment.createdAt,
			PostID:          uint(comment.postId),
			ParentCommentID: comment.parentCommentId,
			RepliesIDs:      commentsIds,
		}
	}

	return comments, nil
}

func (s *Storage) ToggleComments(ctx context.Context, postId uint, userId uint) (bool, error) {
	post, ok := s.posts.Load(uint64(postId))
	if !ok {
		return false, server.ErrPostNotFound
	}

	if post.authorId != uint64(userId) {
		return false, server.ErrUnauthorized
	}

	post.commentsAvailable = !post.commentsAvailable
	s.posts.Store(uint64(postId), post)

	return true, nil
}

func (s *Storage) GetPostsFromUser(ctx context.Context, userId uint) ([]*models.Post, error) {
	posts := make([]*models.Post, 0)

	s.posts.Range(func(id uint64, p *Post) bool {
		if p.authorId == uint64(userId) {
			commentsIds := make([]uint, 0)
			s.comments.Range(func(id uint64, c *Comment) bool {
				if c.postId == id {
					commentsIds = append(commentsIds, uint(c.id))
				}
				return true
			})

			posts = append(posts, &models.Post{
				ID:          uint(p.id),
				Title:       p.title,
				CreatedAt:   p.createdAt,
				Content:     p.content,
				AuthorID:    uint(p.authorId),
				CommentsIDs: commentsIds,
			})
		}
		return true
	})

	return posts, nil
}

func (s *Storage) GetReplies(ctx context.Context, commentId uint) ([]*models.Comment, error) {
	replies := make([]*models.Comment, 0)
	s.comments.Range(func(id uint64, c *Comment) bool {
		if c.parentCommentId != nil && *c.parentCommentId == commentId {
			commentsIds := make([]uint, 0)
			s.comments.Range(func(id uint64, c *Comment) bool {
				if c.postId == id {
					commentsIds = append(commentsIds, uint(c.id))
				}
				return true
			})
			replies = append(replies, &models.Comment{
				ID:              uint(c.id),
				Content:         c.content,
				AuthorID:        uint(c.authorId),
				CreatedAt:       c.createdAt,
				PostID:          uint(c.postId),
				ParentCommentID: c.parentCommentId,
				RepliesIDs:      commentsIds,
			})
		}
		return true
	})

	return replies, nil
}

func (s *Storage) GetCommentsForPost(ctx context.Context, postId uint) ([]*models.Comment, error) {
	//TODO implement me
	panic("implement me")
}
