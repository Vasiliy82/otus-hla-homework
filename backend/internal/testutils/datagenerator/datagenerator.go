package datagenerator

import "github.com/Vasiliy82/otus-hla-homework/domain"

type DataGenerator interface {
	CreateUser(*domain.User) (domain.UserKey, error)
	CreatePost(userId domain.UserKey, message domain.PostMessage) (domain.PostKey, error)
	AddFriend(userId, friendId domain.UserKey) error
}
