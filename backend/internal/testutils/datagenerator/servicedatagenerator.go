package datagenerator

import "github.com/Vasiliy82/otus-hla-homework/domain"

type ServiceDataGenerator struct {
	userService domain.UserService
	postService domain.PostService
}

func NewServiceDataGenerator(us domain.UserService, ps domain.PostService) *ServiceDataGenerator {
	return &ServiceDataGenerator{
		userService: us,
		postService: ps,
	}
}

func (s *ServiceDataGenerator) CreateUser(user *domain.User) (domain.UserKey, error) {
	return s.userService.RegisterUser(user)
}

func (s *ServiceDataGenerator) CreatePost(userId domain.UserKey, message domain.PostMessage) (domain.PostKey, error) {
	return s.postService.Create(userId, message)
}

func (s *ServiceDataGenerator) AddFriend(userId, friendId domain.UserKey) error {
	return s.userService.AddFriend(userId, friendId)
}
