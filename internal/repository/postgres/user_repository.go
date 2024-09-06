package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(user domain.User) (string, error) {
	var userId string
	err := r.db.QueryRow("INSERT INTO users (first_name, second_name, birthdate, biography, city, username, password_hash) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		user.FirstName, user.SecondName, user.Birthdate, user.Biography, user.City, user.Username, user.PasswordHash).Scan(&userId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // duplicate key value violates unique constraint
				return "", domain.ErrConflict
			}
		}
		return "", err
	}
	return userId, nil
}

func (r *UserRepository) GetUserByID(id string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, first_name, second_name, birthdate, biography, city, username, created_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City, &user.Username, &user.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) CheckUserPasswordHash(username string, passwordHash string) (string, error) {
	var id string
	query := "SELECT id FROM users WHERE username = $1 AND password_hash = $2"
	err := r.db.QueryRow(query, username, passwordHash).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("auth error: %v", err)
	}
	return id, nil
}
