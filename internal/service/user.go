package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/lib/pq"
)

//go:generate mockery --name UserRepository
type UserRepository interface {
	RegisterUser(user domain.User) (string, error)
	GetByID(id string) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
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
	var id string
	var err error

	if id, err = s.userRepo.RegisterUser(user); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // duplicate key value violates unique constraint
				return "", apperrors.NewConflictError("Login already used")
			}
		}
		// Если ошибка не является *pq.Error, оборачиваем её в InternalServerError
		return "", apperrors.NewInternalServerError("UserService.RegisterUser, s.userRepo.RegisterUser returned unknown error", err)
	}
	return id, nil
}

func (s *UserService) GetById(id string) (domain.User, error) {
	var user domain.User
	var err error
	if user, err = s.userRepo.GetByID(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, apperrors.NewNotFoundError("User not found")
		}
		return domain.User{}, apperrors.NewInternalServerError("UserService.GetById: s.userRepo.GetByID returned unknown error", err)
	}
	return user, nil

}

func (s *UserService) Login(username, password string) (string, error) {
	// Проверка пароля
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperrors.NewNotFoundError("User not found")
		}
		return "", apperrors.NewInternalServerError("UserService.Login: s.userRepo.GetByUserName returned unknown error", err)
	}
	if !user.CheckPassword(password) {
		return "", apperrors.NewUnauthorizedError("Wrong password")
	}

	// Генерация токена
	token := generateToken(username)

	// Установка срока действия токена (например, на 24 часа)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Сохранение токена в таблицу сессий
	err = s.sessionRepo.CreateSession(user.ID, token, expiresAt)
	if err != nil {
		return "", apperrors.NewInternalServerError("UserSevice.Login: s.sessionRepo.CreateSession returned unknown error", err)
	}

	return token, nil
}

func generateToken(username string) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%s:%d", username, timestamp)
	return domain.HashPassword(data)
}
