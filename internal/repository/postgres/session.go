package postgres

import (
	"fmt"
	"time"
)

func (r *UserRepository) CreateSession(userID, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	return nil
}
