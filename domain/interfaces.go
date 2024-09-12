package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name UserRepository
type UserRepository interface {
	RegisterUser(user User) (string, error)
	GetByID(id string) (User, error)
	GetByUsername(username string) (User, error)
}

//go:generate mockery --name SessionRepository
type SessionRepository interface {
	CreateSession(userID string, expiresAt time.Time) (string, error)
}

//go:generate mockery --name UserService
type UserService interface {
	RegisterUser(user User) (string, error)
	GetById(id string) (User, error)
	Login(username, password string) (string, string, error)
}

//go:generate mockery --name JWTService
type JWTService interface {
	GenerateToken(userID string, permissions []string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
