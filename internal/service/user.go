package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
)

//go:generate mockery --name UserRepository
type UserRepository interface {
	RegisterUser(user domain.User) (string, error)
	GetUserByID(id string) (domain.User, error)
	CheckUserPasswordHash(username string, passwordHash string) (string, error)
}

//go:generate mockery --name SessionRepository
type SessionRepository interface {
	CreateSession(userID, token string, expiresAt time.Time) error
}

type UserService struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewUserService(ur UserRepository, sr SessionRepository) *UserService {
	return &UserService{
		userRepo:    ur,
		sessionRepo: sr,
	}
}

func (s *UserService) RegisterUser(user domain.User) (string, error) {
	return s.userRepo.RegisterUser(user)
}

func (s *UserService) GetById(id string) (domain.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) Login(username, password string) (string, error) {
	passwordHash := domain.HashPassword(password)

	// Проверка пароля
	id, err := s.userRepo.CheckUserPasswordHash(username, passwordHash)
	if err != nil {
		return "", errors.New("Auth error")
	}
	if id == "" {
		return "", errors.New("Auth error")
	}

	// Генерация токена
	token := generateToken(username)

	// Установка срока действия токена (например, на 24 часа)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Сохранение токена в таблицу сессий
	err = s.sessionRepo.CreateSession(id, token, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return token, nil
}

func generateToken(username string) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%s:%d", username, timestamp)
	return domain.HashPassword(data)
}
