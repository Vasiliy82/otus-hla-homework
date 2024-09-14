package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/domain/mocks"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest"
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
	mockUserService := mocks.NewUserService(t)
	mockUserService.On("GetById", "123").Return(&domain.User{ID: "123", Username: "testuser"}, nil)

	// Создаем обработчик
	h := rest.NewUserHandler(mockUserService)

	// Вызываем метод Get
	err := h.Get(c)

	// Проверяем успешный ответ
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "testuser")

	mockUserService.AssertExpectations(t)
}
