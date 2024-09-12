package rest

import (
	"errors"
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
	"github.com/Vasiliy82/otus-hla-homework/internal/mappers"
	log "github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/internal/services/user"
	"github.com/Vasiliy82/otus-hla-homework/internal/validators"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService domain.UserService
}

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
			return c.JSON(http.StatusBadRequest, appverr)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if user, err = mappers.ToUser(userReq); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userId, err := h.userService.RegisterUser(user)
	if err != nil {
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

	user_id, token, err := h.userService.Login(req.Username, req.Password)

	if err != nil {
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"user_id": user_id, "token": token})
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
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}
