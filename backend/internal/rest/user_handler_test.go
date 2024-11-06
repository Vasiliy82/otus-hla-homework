package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/domain/mocks"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Get_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	// Мокаем UserService
	snService := mocks.NewSocialNetworkService(t)
	snService.On("GetUser", "123").Return(&domain.User{ID: "123", Username: "testuser"}, nil)

	// Создаем обработчик
	h := rest.NewSocialNetworkHandler(snService, &config.APIConfig{})

	// Вызываем метод Get
	err := h.GetUser(c)

	// Проверяем успешный ответ
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "testuser")

	snService.AssertExpectations(t)
}

// Тест на невалидный ID
func TestUserHandler_Get_InvalidID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")

	// Создаем мок UserService
	snService := mocks.NewSocialNetworkService(t)

	// Создаем обработчик
	h := rest.NewSocialNetworkHandler(snService, &config.APIConfig{})

	// Вызываем метод Get
	err := h.GetUser(c)

	// Проверяем ответ на невалидный ID
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	snService.AssertExpectations(t)
}

func TestUserHandler_Get_UserNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	// Мокаем UserService
	snService := mocks.NewSocialNetworkService(t)
	snService.On("GetUser", "999").Return(nil, apperrors.NewNotFoundError("User not found"))

	// Создаем обработчик
	h := rest.NewSocialNetworkHandler(snService, &config.APIConfig{})

	// Вызываем метод Get
	err := h.GetUser(c)

	// Проверяем ответ на несуществующего пользователя
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "User not found")

	snService.AssertExpectations(t)
}

// Тест на успешный Logout
func TestUserHandler_Logout_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Мокаем токен и добавляем его в контекст
	token := &jwt.Token{
		Claims: &domain.UserClaims{
			Permissions: []domain.Permission{domain.PermissionUserGet},
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "12345",
			},
		},
	}
	c.Set("token", token)

	// Мокаем UserService
	snService := mocks.NewSocialNetworkService(t)
	snService.On("Logout", token).Return(nil)

	// Создаем обработчик
	h := rest.NewSocialNetworkHandler(snService, &config.APIConfig{})

	// Вызываем метод Logout
	err := h.Logout(c)

	// Проверяем успешный ответ
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "null")

	snService.AssertExpectations(t)
}

// Тест на отсутствие токена при Logout
func TestUserHandler_Logout_MissingToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Мокаем UserService
	snService := mocks.NewSocialNetworkService(t)

	// Создаем обработчик
	h := rest.NewSocialNetworkHandler(snService, &config.APIConfig{})

	// Вызываем метод Logout без токена
	err := h.Logout(c)

	// Проверяем ответ на отсутствие токена
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	snService.AssertExpectations(t)
}

// Тест на ошибку при Logout
func TestUserHandler_Logout_Failed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Мокаем токен и добавляем его в контекст
	token := &jwt.Token{
		Claims: &domain.UserClaims{
			Permissions: []domain.Permission{"read"},
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "12345",
			},
		},
	}
	c.Set("token", token)

	// Мокаем UserService с ошибкой при Logout
	snService := mocks.NewSocialNetworkService(t)
	snService.On("Logout", token).Return(errors.New("logout error"))

	// Создаем обработчик
	h := rest.NewSocialNetworkHandler(snService, &config.APIConfig{})

	// Вызываем метод Logout
	err := h.Logout(c)

	// Проверяем ответ на ошибку
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	snService.AssertExpectations(t)
}
