package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/tauadam/reading_list-bot/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
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
	query := `SELECT articles FROM articles WHERE user_name = ? OR ORDER BY random() LIMIT 1`

	var res string
	err := s.db.QueryRowContext(ctx, query, userName).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
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
