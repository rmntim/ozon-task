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

func (s *Storage) CreateUser(ctx context.Context, username string, email string, password string) (*model.User, error) {
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

func (s *Storage) CreatePost(ctx context.Context, title string, content string, authorId uint) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CreateComment(ctx context.Context, content string, id uint, id2 uint, id3 *uint) (*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetUserById(ctx context.Context, id uint) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetUsers(ctx context.Context) ([]*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetPostById(ctx context.Context, id uint) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetPosts(ctx context.Context) ([]*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetCommentById(ctx context.Context, id uint) (*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetComments(ctx context.Context) ([]*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}
