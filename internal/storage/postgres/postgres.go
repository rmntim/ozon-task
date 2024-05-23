package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rmntim/ozon-task/graph/model"
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

func (s *Storage) GetPostById(ctx context.Context, id int) (*model.Post, error) {
	const op = "storage.postgres.GetPostById"

	stmt, err := s.db.PrepareContext(ctx, `
					SELECT p.id,
       					   p.title,
       					   u.id,
       					   u.username,
       					   u.email,
       					   p.created_at,
       					   p.content,
       					   p.allow_comments
					FROM posts p
         				JOIN users u ON p.creator_id = u.id
					WHERE p.id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var post model.Post
	// Imagine using ptrs as default
	post.Creator = &model.User{}
	if err := rows.Scan(&post.ID, &post.Title,
		&post.Creator.ID, &post.Creator.Name, &post.Creator.Email,
		&post.CreatedAt, &post.Content, &post.AllowComments); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	post.Comments = make([]*model.Comment, 0)

	return &post, nil
}

func (s *Storage) CreatePost(ctx context.Context, post model.PostInput) (*model.Post, error) {
	const op = "storage.postgres.CreatePost"

	var id int
	if err := s.db.QueryRowContext(ctx, `
		INSERT INTO posts (title, creator_id, content, allow_comments)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, post.Title, post.CreatorID, post.Content, post.AllowComments).Scan(&id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.GetPostById(ctx, id)
}
