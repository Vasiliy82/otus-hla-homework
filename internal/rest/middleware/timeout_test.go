package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/internal/rest/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetRequestContextWithTimeout_Success(t *testing.T) {
	// Создаем новый echo инстанс
	e := echo.New()

	// Устанавливаем длительность таймаута
	timeoutDuration := 100 * time.Millisecond

	// Создаем middleware с таймаутом
	middlewareFunc := middleware.SetRequestContextWithTimeout(timeoutDuration)

	// Мокаем обработчик, который будет проверять наличие контекста
	handler := func(c echo.Context) error {
		ctx := c.Request().Context()

		// Проверяем, что контекст еще активен (не истек)
		deadline, ok := ctx.Deadline()
		require.True(t, ok, "контекст должен иметь таймаут")
		require.WithinDuration(t, time.Now().Add(timeoutDuration), deadline, time.Millisecond*10)

		return c.String(http.StatusOK, "success")
	}

	// Применяем middleware к нашему обработчику
	wrappedHandler := middlewareFunc(handler)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Создаем echo-контекст
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	err := wrappedHandler(c)

	// Проверяем, что ошибок не возникло
	assert.NoError(t, err)

	// Проверяем, что обработчик возвращает правильный статус и ответ
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "success", rec.Body.String())
}

func TestSetRequestContextWithTimeout_Timeout(t *testing.T) {
	// Создаем новый echo инстанс
	e := echo.New()

	// Устанавливаем короткий таймаут
	timeoutDuration := 50 * time.Millisecond

	// Создаем middleware с таймаутом
	middlewareFunc := middleware.SetRequestContextWithTimeout(timeoutDuration)

	// Мокаем обработчик, который имитирует длительную операцию
	handler := func(c echo.Context) error {
		ctx := c.Request().Context()

		// Имитируем длительную операцию
		select {
		case <-ctx.Done():
			// Если контекст завершен по таймауту, возвращаем ошибку
			return c.String(http.StatusGatewayTimeout, "request timed out")
		case <-time.After(100 * time.Millisecond):
			// Если операция завершилась раньше таймаута, возвращаем успех
			return c.String(http.StatusOK, "success")
		}
	}

	// Применяем middleware к нашему обработчику
	wrappedHandler := middlewareFunc(handler)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Создаем echo-контекст
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	err := wrappedHandler(c)

	// Проверяем, что ошибок не возникло
	assert.NoError(t, err)

	// Проверяем, что обработчик возвращает статус 504 (таймаут)
	assert.Equal(t, http.StatusGatewayTimeout, rec.Code)
	assert.Equal(t, "request timed out", rec.Body.String())
}
