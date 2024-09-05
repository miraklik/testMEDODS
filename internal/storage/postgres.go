package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

type Storage struct {
	db *pgx.Conn
}

func NewStorage(connStr string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &Storage{db: conn}, nil
}

func (s *Storage) StoreRefreshToken(userID, hashedToken string) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO refresh_tokens (user_id, token, created_at) VALUES ($1, $2, $3)",
		userID, hashedToken, time.Now())
	return err
}

func (s *Storage) GetRefreshToken(userID string) (string, error) {
	var hashedToken string
	err := s.db.QueryRow(context.Background(),
		"SELECT token FROM refresh_tokens WHERE user_id=$1", userID).Scan(&hashedToken)
	if err != nil {
		return "", err
	}
	return hashedToken, nil
}

func (s *Storage) DeleteRefreshToken(userID string) error {
	_, err := s.db.Exec(context.Background(), "DELETE FROM refresh_tokens WHERE user_id=$1", userID)
	return err
}
