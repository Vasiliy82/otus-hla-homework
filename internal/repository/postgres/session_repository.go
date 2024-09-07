package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(userID, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("SessionRepository.CreateSession: r.db.Exec returned error %w", err)
	}
	return nil
}
