package rest

import (
	"log"
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
	"github.com/Vasiliy82/otus-hla-homework/internal/mappers"
	"github.com/Vasiliy82/otus-hla-homework/internal/service"
	"github.com/Vasiliy82/otus-hla-homework/internal/validators"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	var userReq dto.RegisterUserRequest
	var user domain.User
	var err error

	log.Println("RegisterUser")

	if err = c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err = validators.ValidateRegisterUserRequest(userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if user, err = mappers.ToUser(userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userId, err := h.userService.RegisterUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"user_id": userId})
}

func (h *UserHandler) Login(c echo.Context) error {

	var req dto.LoginRequest
	var err error

	log.Println("Login")

	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err = validators.ValidateLoginRequest(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) Get(c echo.Context) error {
	var err error

	id := c.Param("id")

	if err = validators.ValidateUserId(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userService.GetById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}
