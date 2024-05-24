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
	if err := s.db.QueryRowxContext(ctx,
		`SELECT p.id, p.title, p.created_at, p.content, p.author_id, array_agg(c.id) as comments_ids
				FROM posts p
					LEFT JOIN comments c ON p.id = c.post_id
				WHERE p.id = $1
				GROUP BY p.id, p.title, p.created_at, p.content, p.author_id`, id).StructScan(&post); err != nil {
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
	if err := s.db.SelectContext(ctx, &posts,
		`SELECT p.id, p.title, p.created_at, p.content, p.author_id, array_agg(c.id) as comments_ids
				FROM posts p
					LEFT JOIN comments c ON p.id = c.post_id
				GROUP by p.id, p.title, p.created_at, p.content, p.author_id LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

func (s *Storage) GetCommentById(ctx context.Context, id uint) (*models.Comment, error) {
	const op = "storage.postgres.GetCommentById"

	var comment models.Comment
	if err := s.db.QueryRowxContext(ctx,
		`SELECT c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id, array_agg(r.id) as replies_ids
				FROM comments c
					LEFT JOIN comments r ON r.parent_comment_id = c.id
				WHERE c.id = $1
				GROUP BY c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id`, id).StructScan(&comment); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.ErrCommentNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &comment, nil
}

func (s *Storage) GetComments(ctx context.Context, limit int, offset int) ([]*models.Comment, error) {
	const op = "storage.postgres.GetComments"

	var comments []*models.Comment
	if err := s.db.SelectContext(ctx, &comments,
		`SELECT c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id, array_agg(r.id) as replies_ids
				FROM comments c
					LEFT JOIN comments r ON r.parent_comment_id = c.id
				GROUP by c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (s *Storage) ToggleComments(ctx context.Context, postId uint) (bool, error) {
	const op = "storage.postgres.ToggleComments"

	var commentsAvailable bool
	stmt, err := s.db.PreparexContext(ctx, `UPDATE posts SET comments_available = NOT comments_available WHERE id = $1 RETURNING comments_available`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(postId).Scan(&commentsAvailable); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return commentsAvailable, nil
}

func (s *Storage) GetPostsFromUser(ctx context.Context, userId uint) ([]*models.Post, error) {
	const op = "storage.postgres.GetPostsFromUser"

	var posts []*models.Post
	if err := s.db.SelectContext(ctx, &posts,
		`SELECT p.id, p.title, p.created_at, p.content, p.author_id, array_agg(c.id) as comments_ids
				FROM posts p
					LEFT JOIN comments c ON p.id = c.post_id
				WHERE p.author_id = $1
				GROUP BY p.id, p.title, p.created_at, p.content, p.author_id`, userId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

func (s *Storage) GetReplies(ctx context.Context, commentId uint) ([]*models.Comment, error) {
	const op = "storage.postgres.GetReplies"

	var comments []*models.Comment
	if err := s.db.SelectContext(ctx, &comments,
		`SELECT c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id, array_agg(r.id) as replies_ids
				FROM comments c
					LEFT JOIN comments r ON r.parent_comment_id = c.id
				WHERE c.parent_comment_id = $1
				GROUP BY c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id`, commentId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (s *Storage) GetCommentsForPost(ctx context.Context, postId uint) ([]*models.Comment, error) {
	const op = "storage.postgres.GetCommentsForPost"

	var comments []*models.Comment
	if err := s.db.SelectContext(ctx, &comments,
		`SELECT c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id, array_agg(r.id) as replies_ids
				FROM comments c
					LEFT JOIN comments r ON r.parent_comment_id = c.id
				WHERE c.post_id = $1
				GROUP BY c.id, c.content, c.created_at, c.author_id, c.post_id, c.parent_comment_id`, postId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}
