package sqlite

import (
	"context"
	"database/sql"
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

func (s *Storage) Save(ctx context.Context, a *storage.Article) error {
	query := `INSERT INTO articles (url, user_name) VALUES (?,?)`

	_, err := s.db.ExecContext(ctx, query, a.URL, a.UserName)
	if err != nil {
		return fmt.Errorf("failed to save article: %w", err)
	}

	return nil
}
