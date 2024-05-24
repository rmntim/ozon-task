package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rmntim/ozon-task/internal/models"
	"github.com/rmntim/ozon-task/internal/server"
)

type Storage struct {
	db *sqlx.DB
}

// New creates new postgres storage instance and pings it to check connection.
func New(username, password, address, database string) (*Storage, error) {
	const op = "storage.postgres.New"

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, address, database)

	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Migrate runs migrations.
func (s *Storage) Migrate() error {
	const op = "storage.postgres.Migrate"

	driver, err := postgres.WithInstance(s.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/postgres",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CreateUser(ctx context.Context, username string, email string, password string) (*models.User, error) {
	const op = "storage.postgres.CreateUser"

	// FIXME: i dont want to bother with password hashing, lets just imagine it works
	passwordHash := []byte(password + "_hashed")

	stmt, err := s.db.PreparexContext(ctx, "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(username, email, passwordHash).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.GetUserById(ctx, id)
}

func (s *Storage) CreatePost(ctx context.Context, title string, content string, authorId uint) (*models.Post, error) {
	const op = "storage.postgres.CreatePost"

	stmt, err := s.db.PreparexContext(ctx, "INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(title, content, authorId).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.GetPostById(ctx, id)
}

func (s *Storage) CreateComment(ctx context.Context, content string, authorId uint, postId uint, parentCommentId *uint) (*models.Comment, error) {
	const op = "storage.postgres.CreateComment"

	stmt, err := s.db.PreparexContext(ctx, "INSERT INTO comments (content, author_id, post_id, parent_comment_id) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(content, authorId, postId, parentCommentId).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.GetCommentById(ctx, id)
}

func (s *Storage) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	const op = "storage.postgres.GetUserById"

	var user models.User
	if err := s.db.QueryRowxContext(ctx, "SELECT id, username, email FROM users WHERE id = $1", id).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) GetUsers(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	const op = "storage.postgres.GetUsers"

	var users []*models.User
	if err := s.db.SelectContext(ctx, &users, "SELECT id, username, email FROM users LIMIT $1 OFFSET $2", limit, offset); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *Storage) GetPostById(ctx context.Context, id uint) (*models.Post, error) {
	const op = "storage.postgres.GetPostById"

	var post models.Post
	post.Author = &models.User{}
	if err := s.db.QueryRowxContext(ctx, "SELECT p.id, title, created_at, content, u.id, username, email FROM posts p JOIN users u ON p.author_id = u.id WHERE p.id = $1", id).
		Scan(&post.ID, &post.Title, &post.CreatedAt, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.ErrPostNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &post, nil
}

func (s *Storage) GetPosts(ctx context.Context, limit int, offset int) ([]*models.Post, error) {
	const op = "storage.postgres.GetPosts"

	var posts []*models.Post
	rows, err := s.db.QueryxContext(ctx,
		"SELECT p.id, title, created_at, content, u.id, username, email FROM posts p JOIN users u ON p.author_id = u.id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		post.Author = &models.User{}
		if err := rows.Scan(&post.ID, &post.Title, &post.CreatedAt, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *Storage) GetCommentById(ctx context.Context, id uint) (*models.Comment, error) {
	const op = "storage.postgres.GetCommentById"

	var comment models.Comment
	var authorId uint
	var postId uint
	var parentCommentId *uint

	if err := s.db.QueryRowxContext(ctx, "SELECT id, content, author_id, post_id, parent_comment_id FROM comments WHERE id = $1", id).
		Scan(&comment.ID, &comment.Content, &authorId, &postId, &parentCommentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.ErrCommentNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var err error
	comment.Author, err = s.GetUserById(ctx, authorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	comment.Post, err = s.GetPostById(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if parentCommentId != nil {
		comment.ParentComment, err = s.GetCommentById(ctx, *parentCommentId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &comment, nil
}

func (s *Storage) GetComments(ctx context.Context, limit int, offset int) ([]*models.Comment, error) {
	const op = "storage.postgres.GetComments"

	var comments []*models.Comment
	var authorIds []uint
	var postIds []uint
	var parentCommentIds []*uint
	rows, err := s.db.QueryxContext(ctx, "SELECT id, content, created_at, author_id, post_id, parent_comment_id FROM comments LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		var authorId uint
		var postId uint
		var parentCommentId *uint
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &authorId, &postId, &parentCommentId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		comments = append(comments, &comment)
		authorIds = append(authorIds, authorId)
		postIds = append(postIds, postId)
		parentCommentIds = append(parentCommentIds, parentCommentId)
	}

	// FIXME: n+1 my love
	for i := 0; i < len(comments); i++ {
		comments[i].Author, err = s.GetUserById(ctx, authorIds[i])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments[i].Post, err = s.GetPostById(ctx, postIds[i])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if parentCommentIds[i] != nil {
			comments[i].ParentComment, err = s.GetCommentById(ctx, *parentCommentIds[i])
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	return comments, nil
}
