package rest

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
	"github.com/Vasiliy82/otus-hla-homework/internal/mappers"
	log "github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/internal/services/user"
	"github.com/Vasiliy82/otus-hla-homework/internal/validators"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService domain.UserService
}

// Регулярное выражение для проверки только букв
var validNameRegex = regexp.MustCompile(`^[\p{L}]+$`) // \p{L} соответствует любому юникодовскому символу, который является буквой

func NewUserHandler(userService domain.UserService) user.UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) RegisterUser(c echo.Context) error {
	var userReq dto.RegisterUserRequest
	var user domain.User
	var err error
	log.Logger().Debug("UserHandler.RegisterUser")

	if err = c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err = validators.ValidateRegisterUserRequest(userReq); err != nil {
		var appverr *apperrors.ValidationError
		if errors.As(err, &appverr) {
			log.Logger().Warnw("userHandler.RegisterUser: validators.ValidateRegisterUserRequest returned ValidationError", "err", err)
			return c.JSON(http.StatusBadRequest, appverr)
		}
		log.Logger().Errorw("userHandler.RegisterUser: validators.ValidateRegisterUserRequest returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if user, err = mappers.ToUser(userReq); err != nil {
		log.Logger().Errorw("userHandler.RegisterUser: mappers.ToUser returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userId, err := h.userService.RegisterUser(&user)
	if err != nil {
		log.Logger().Errorw("userHandler.RegisterUser: h.userService.RegisterUser returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, apperr)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := dto.LoginResponse{ID: userId, Token: ""}

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) Login(c echo.Context) error {

	var req dto.LoginRequest
	var err error

	log.Logger().Debug("UserHandler.Login")

	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err = validators.ValidateLoginRequest(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.userService.Login(req.Username, req.Password)

	if err != nil {
		log.Logger().Errorw("userHandler.Login: h.userService.Login returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": string(token)})
}

func (h *userHandler) Get(c echo.Context) error {
	var err error

	log.Logger().Debug("UserHandler.Get")

	id := c.Param("id")

	if err = validators.ValidateUserId(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userService.GetById(id)
	if err != nil {
		log.Logger().Errorw("userHandler.Get: h.userService.GetById returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": apperr.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) Search(c echo.Context) error {
	// Извлечение query параметров first_name и last_name
	firstName := c.QueryParam("first_name")
	lastName := c.QueryParam("last_name")

	// Валидация параметров
	if !isValidName(firstName) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат имени"})
	}
	if !isValidName(lastName) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат фамилии"})
	}

	users, err := h.userService.Search(firstName, lastName)
	if err != nil {
		log.Logger().Errorw("userHandler.Search: h.userService.Search returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": apperr.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *userHandler) AddFriend(c echo.Context) error {
	log.Logger().Debug("UserHandler.AddFriend")

	var err error

	// Извлечение идентификатора будущего друга из URL
	friend_id := c.Param("friend_id")

	if err = validators.ValidateUserId(friend_id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Извлечение своего идентификатора из контекста
	// Здесь наверное нужен рефакторинг, сделать приватный метод со всеми проверками,
	// т.к. в будущем контекст пользователя понадобится практически везде
	claims, ok := c.Get("claims").(*domain.UserClaims)
	if !ok {
		// Теоретически, такого не должно случиться, т.к. токен проверяется в Middleware
		return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid token"))
	}

	my_id := claims.Subject

	// Эту проверку имеет смысл вынести в middleware, чтобы валидировать в одном месте
	if err = validators.ValidateUserId(my_id); err != nil {
		log.Logger().Errorw("userHandler.AddFriend: validators.ValidateUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}

	if err = h.userService.AddFriend(my_id, friend_id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrFriendAlreadyExists) {
			log.Logger().Warnw("userHandler.AddFriend: h.userService.AddFriend returned ErrUserNotFound", "err", err)
			return c.JSON(http.StatusBadRequest, apperrors.NewBadRequestError(err.Error()))
		}
		log.Logger().Errorw("userHandler.AddFriend: h.userService.AddFriend returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *userHandler) RemoveFriend(c echo.Context) error {
	log.Logger().Debug("UserHandler.RemoveFriend")

	var err error

	// Извлечение идентификатора будущего друга из URL
	friend_id := c.Param("friend_id")

	if err = validators.ValidateUserId(friend_id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Извлечение своего идентификатора из контекста
	// Здесь наверное нужен рефакторинг, сделать приватный метод со всеми проверками,
	// т.к. в будущем контекст пользователя понадобится практически везде
	claims, ok := c.Get("claims").(*domain.UserClaims)
	if !ok {
		// Теоретически, такого не должно случиться, т.к. токен проверяется в Middleware
		return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid token"))
	}

	my_id := claims.Subject

	// Эту проверку имеет смысл вынести в middleware, чтобы валидировать в одном месте
	if err = validators.ValidateUserId(my_id); err != nil {
		log.Logger().Errorw("userHandler.RemoveFriend: validators.ValidateUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}

	if err = h.userService.RemoveFriend(my_id, friend_id); err != nil {
		if errors.Is(err, domain.ErrFriendNotFound) {
			log.Logger().Warnw("userHandler.RemoveFriend: h.userService.RemoveFriend returned ErrFriendNotFound", "err", err)
			return c.JSON(http.StatusBadRequest, apperrors.NewBadRequestError(err.Error()))
		}
		log.Logger().Errorw("userHandler.RemoveFriend: h.userService.RemoveFriend returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *userHandler) Logout(c echo.Context) error {
	log.Logger().Debug("UserHandler.Logout")
	// Извлекаем токен из контекста
	token, ok := c.Get("token").(*jwt.Token)
	if !ok {
		// Теоретически такого не может случиться, т.к. токен проверяется в Middleware
		return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid token"))
	}

	if err := h.userService.Logout(token); err != nil {
		log.Logger().Errorw("userHandler.Logout: h.userService.Logout returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}

	return c.JSON(http.StatusOK, nil)
}

// Функция для валидации имени
func isValidName(name string) bool {
	// Проверяем, что строка содержит только буквы
	return validNameRegex.MatchString(name)
}
