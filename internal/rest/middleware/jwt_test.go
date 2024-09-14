package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/domain/mocks"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Создаем валидный токен с кастомными claims
func createValidToken() *jwt.Token {
	return &jwt.Token{
		Claims: &domain.UserClaims{
			Permissions: []domain.Permission{domain.PermissionUserGet},
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "12345", // ID пользователя
			},
		},
		Method: jwt.SigningMethodHS256,
	}
}

func TestJWTMiddleware_MissingAuthorizationHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Создаем мок JWTService через mocks.NewJWTService
	mockJWTService := mocks.NewJWTService(t)

	// Применяем middleware
	middlewareFunc := middleware.JWTMiddleware(mockJWTService)
	handler := middlewareFunc(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Вызываем middleware
	err := handler(c)

	// Проверяем, что вернулся статус 401 Unauthorized
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidHeaderFormat") // Некорректный формат заголовка
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Создаем мок JWTService через mocks.NewJWTService
	mockJWTService := mocks.NewJWTService(t)

	// Применяем middleware
	middlewareFunc := middleware.JWTMiddleware(mockJWTService)
	handler := middlewareFunc(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Вызываем middleware
	err := handler(c)

	// Проверяем, что вернулся статус 401 Unauthorized
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidToken") // Невалидный токен
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Создаем мок JWTService через mocks.NewJWTService
	mockJWTService := mocks.NewJWTService(t)
	mockJWTService.On("ValidateToken", domain.TokenString("invalidToken")).Return(nil, errors.New("invalid token"))

	// Применяем middleware
	middlewareFunc := middleware.JWTMiddleware(mockJWTService)
	handler := middlewareFunc(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Вызываем middleware
	err := handler(c)

	// Проверяем, что вернулся статус 401 Unauthorized
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer validToken") // Валидный токен
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Создаем мок JWTService через mocks.NewJWTService
	mockJWTService := mocks.NewJWTService(t)
	validToken := createValidToken() // Пример валидного токена с кастомными claims
	mockJWTService.On("ValidateToken", domain.TokenString("validToken")).Return(validToken, nil)

	// Применяем middleware
	middlewareFunc := middleware.JWTMiddleware(mockJWTService)
	handler := middlewareFunc(func(c echo.Context) error {
		// Проверяем, что токен сохранен в контексте
		token, ok := c.Get("token").(*jwt.Token)
		assert.True(t, ok)
		assert.Equal(t, validToken, token)

		// Проверяем кастомные claims
		claims, ok := token.Claims.(*domain.UserClaims)
		assert.True(t, ok)
		assert.Equal(t, "12345", claims.Subject)
		assert.Equal(t, []domain.Permission{domain.PermissionUserGet}, claims.Permissions)

		return c.String(http.StatusOK, "success")
	})

	// Вызываем middleware
	err := handler(c)

	// Проверяем, что ошибок нет и вернулся успешный статус
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
