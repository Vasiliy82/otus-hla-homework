package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenString string

type Permission string

type Token struct {
	Serial      int64
	Subject     string
	Expire      time.Time
	Permissions []Permission
	JWTToken    *jwt.Token
}

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

type BlacklistRepository interface {
	NewSerial() (int64, error)
	AddToBlacklist(serial int64, expireDate time.Time) error
	IsBlacklisted(serial int64) (bool, error)
}

//go:generate mockery --name UserService
type UserService interface {
	RegisterUser(user User) (string, error)
	GetById(id string) (User, error)
	Login(username, password string) (TokenString, error)
}

//go:generate mockery --name JWTService
type JWTService interface {
	GenerateToken(userID string, permissions []Permission) (TokenString, error)
	ValidateToken(tokenString TokenString) (*Token, error)
	RevokeToken(tokenString TokenString) error
}
