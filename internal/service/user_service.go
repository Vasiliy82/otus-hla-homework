package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/repository/postgres"
)

type UserService struct {
	userRepo *postgres.UserRepository
}

func NewUserService(userRepo *postgres.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
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
	err = s.userRepo.CreateSession(id, token, expiresAt)
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
