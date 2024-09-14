package repository

import (
	"database/sql"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) RegisterUser(user *domain.User) (string, error) {
	var userId string
	err := r.db.QueryRow("INSERT INTO users (first_name, second_name, birthdate, biography, city, username, password_hash) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		user.FirstName, user.SecondName, user.Birthdate, user.Biography, user.City, user.Username, user.PasswordHash).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("userRepository.RegisterUser: r.db.QueryRow returned error %w", err)
	}
	return userId, nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, first_name, second_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &domain.User{}, fmt.Errorf("userRepository.GetByID: r.db.QueryRow returned error %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, first_name, second_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &domain.User{}, fmt.Errorf("userRepository.GetByUsername: r.db.QueryRow returned error %w", err)
	}
	return &user, nil
}
