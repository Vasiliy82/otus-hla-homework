package datagenerator

import "github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"

type DataGenerator interface {
	CreateUser(*domain.User) (domain.UserKey, error)
	CreatePost(userId domain.UserKey, message domain.PostText) (domain.PostKey, error)
	AddFriend(userId, friendId domain.UserKey) error
}
