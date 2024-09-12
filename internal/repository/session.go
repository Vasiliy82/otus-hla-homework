package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
)

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(userID string, expiresAt time.Time) (string, error) {
	token := generateToken(userID)
	_, err := r.db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt)
	if err != nil {
		return "", fmt.Errorf("sessionRepository.CreateSession: r.db.Exec returned error %w", err)
	}
	return token, nil
}

func generateToken(userID string) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%s:%d", userID, timestamp)
	return domain.HashPassword(data)
}
