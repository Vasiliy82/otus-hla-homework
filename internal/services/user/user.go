package user

import (
	"database/sql"
	"errors"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type UserHandler interface {
	RegisterUser(c echo.Context) error
	Login(c echo.Context) error
	Get(c echo.Context) error
	Logout(c echo.Context) error
}

type userService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	jwtService  domain.JWTService
}

func NewUserService(ur domain.UserRepository, sr domain.SessionRepository, jwts domain.JWTService) domain.UserService {
	return &userService{
		userRepo:    ur,
		sessionRepo: sr,
		jwtService:  jwts,
	}
}

func (s *userService) RegisterUser(user domain.User) (string, error) {
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

func (s *userService) GetById(id string) (domain.User, error) {
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

func (s *userService) Login(username, password string) (domain.TokenString, error) {
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

	token, err := s.jwtService.GenerateToken(user.ID, []domain.Permission{domain.PermissionUserGet})
	if err != nil {
		return "", apperrors.NewInternalServerError("UserSevice.Login: s.sessionRepo.CreateSession returned unknown error", err)
	}

	return token, nil
}

func (s *userService) Logout(token *domain.Token) error {

	if err := s.jwtService.RevokeToken(token); err != nil {
		return apperrors.NewInternalServerError("Internal server error", err)
	}
	return nil
}
