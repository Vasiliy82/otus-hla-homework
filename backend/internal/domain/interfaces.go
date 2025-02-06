package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenString string

//go:generate mockery --name UserRepository
type UserRepository interface {
	RegisterUser(user *User) (UserKey, error)
	GetByID(id UserKey) (*User, error)
	GetByUsername(username string) (*User, error)
	Search(firstName, lastName string) ([]*User, error)
	AddFriend(my_id, friend_id UserKey) error
	RemoveFriend(my_id, friend_id UserKey) error
	GetFriendsIds(id UserKey) ([]UserKey, error)
	SetLastActivity(id UserKey) error
	GetUsersActiveSince(period time.Duration) ([]UserKey, error)
}

//go:generate mockery --name PostRepository
type PostRepository interface {
	List(userId UserKey, limit int, lastPostId PostKey) ([]*Post, error)
	Create(userId UserKey, message PostText) (*Post, error)
	Get(postId PostKey) (*Post, error)
	GetPostOwner(postId PostKey) (UserKey, error)
	UpdateMessage(postId PostKey, newMessage PostText) (*Post, error)
	Delete(id PostKey) error
	GetFeed(userId UserKey, limit int) ([]*Post, error)
}

//go:generate mockery --name PostCache
type PostCache interface {
	GetFeed(userId UserKey, limit int) ([]*Post, error)
	UpdateFeed(userId UserKey, posts []*Post) error
}

//go:generate mockery --name BlacklistRepository
type BlacklistRepository interface {
	NewSerial() (string, error)
	AddToBlacklist(serial string, expireDate time.Time) error
	IsBlacklisted(serial string) (bool, error)
}

//go:generate mockery --name SocialNetworkService
type SocialNetworkService interface {
	CreateUser(user *User) (UserKey, error)
	GetUser(id UserKey) (*User, error)
	Search(firstName, lastName string) ([]*User, error)
	Login(username, password string) (TokenString, error)
	AddFriend(my_id, friend_id UserKey) error
	RemoveFriend(my_id, friend_id UserKey) error
	GetFriendsIds(user_id UserKey) ([]UserKey, error)
	Logout(token *jwt.Token) error
	ListPosts(userId UserKey, limit int, lastPostId PostKey) ([]*Post, error)
	CreatePost(userId UserKey, message PostText) (PostKey, error)
	GetPost(userId UserKey, postId PostKey) (*Post, error)
	UpdatePost(userId UserKey, postId PostKey, newMessage PostText) error
	DeletePost(userId UserKey, postId PostKey) error
	GetFeed(userId UserKey) ([]*Post, error)
	SetLastActivity(userId UserKey) error
}

//go:generate mockery --name JWTService
type JWTService interface {
	GenerateToken(userID UserKey, permissions []Permission) (TokenString, error)
	ValidateToken(tokenString TokenString) (*jwt.Token, error)
	RevokeToken(token *jwt.Token) error
	ExtractClaims(token *jwt.Token) (*UserClaims, error)
}
