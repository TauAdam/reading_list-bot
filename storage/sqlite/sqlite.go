package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tauadam/reading_list-bot/storage"
	"os"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("failed to close database file: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

// Init creates the articles table if it doesn't exist. It should be called before any other method
func (s *Storage) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS articles (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// Save saves an article to the database
func (s *Storage) Save(ctx context.Context, a *storage.Article) error {
	query := `INSERT INTO articles (url, user_name) VALUES (?,?)`

	_, err := s.db.ExecContext(ctx, query, a.URL, a.UserName)
	if err != nil {
		return fmt.Errorf("failed to save article: %w", err)
	}

	return nil
}

// PickRandom returns a random article from the database
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Article, error) {
	query := `SELECT url FROM articles WHERE user_name = ? ORDER BY random() LIMIT 1`

	var res string
	err := s.db.QueryRowContext(ctx, query, userName).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrArticleNotFound
		}
		return nil, fmt.Errorf("failed to pick random articles: %w", err)
	}

	return &storage.Article{
		URL:      res,
		UserName: userName,
	}, nil
}

// Remove removes an article from the database
func (s *Storage) Remove(ctx context.Context, article *storage.Article) error {
	query := `DELETE FROM articles WHERE url = ? AND user_name = ?`

	// Prevent SQL injection
	if _, err := s.db.ExecContext(ctx, query, article.URL, article.UserName); err != nil {
		return fmt.Errorf("failed to remove article: %w", err)
	}

	return nil
}

// IsExist checks if an article exists in the database
func (s *Storage) IsExist(ctx context.Context, article *storage.Article) (bool, error) {
	query := `SELECT COUNT(*) FROM articles WHERE url = ? AND user_name = ?`

	var count int

	err := s.db.QueryRowContext(ctx, query, article.URL, article.UserName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if article exists: %w", err)
	}

	return count > 0, nil
}
