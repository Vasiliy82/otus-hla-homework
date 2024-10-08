package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenString string

//go:generate mockery --name UserRepository
type UserRepository interface {
	RegisterUser(user *User) (string, error)
	GetByID(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Search(firstName, lastName string) ([]*User, error)
	AddFriend(my_id, friend_id string) error
	RemoveFriend(my_id, friend_id string) error
}

//go:generate mockery --name BlacklistRepository
type BlacklistRepository interface {
	NewSerial() (string, error)
	AddToBlacklist(serial string, expireDate time.Time) error
	IsBlacklisted(serial string) (bool, error)
}

//go:generate mockery --name UserService
type UserService interface {
	RegisterUser(user *User) (string, error)
	GetById(id string) (*User, error)
	Search(firstName, lastName string) ([]*User, error)
	Login(username, password string) (TokenString, error)
	AddFriend(my_id, friend_id string) error
	RemoveFriend(my_id, friend_id string) error
	Logout(token *jwt.Token) error
}

//go:generate mockery --name JWTService
type JWTService interface {
	GenerateToken(userID string, permissions []Permission) (TokenString, error)
	ValidateToken(tokenString TokenString) (*jwt.Token, error)
	RevokeToken(token *jwt.Token) error
	ExtractClaims(token *jwt.Token) (*UserClaims, error)
}
