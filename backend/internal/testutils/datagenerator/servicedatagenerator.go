package datagenerator

import "github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"

type ServiceDataGenerator struct {
	userService domain.SocialNetworkService
}

func NewServiceDataGenerator(us domain.SocialNetworkService) *ServiceDataGenerator {
	return &ServiceDataGenerator{
		userService: us,
	}
}

func (s *ServiceDataGenerator) CreateUser(user *domain.User) (domain.UserKey, error) {
	return s.userService.CreateUser(user)
}

func (s *ServiceDataGenerator) CreatePost(userId domain.UserKey, message domain.PostText) (domain.PostKey, error) {
	return s.userService.CreatePost(userId, message)
}

func (s *ServiceDataGenerator) AddFriend(userId, friendId domain.UserKey) error {
	return s.userService.AddFriend(userId, friendId)
}
