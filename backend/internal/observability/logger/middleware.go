package logger

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// LoggingMiddleware добавляет x-request-id и x-user-id в контекст логгера
func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		// Извлекаем x-request-id из заголовка или генерируем новый
		requestID := req.Header.Get("x-request-id")
		if requestID == "" {
			requestID = GenerateID()
		}

		// Извлекаем x-user-id из заголовка
		userID := req.Header.Get("x-user-id")

		// Обогащаем логгер
		log := FromContext(req.Context()).With(
			"x-request-id", requestID,
			"x-user-id", userID,
		)

		// Добавляем обогащенный логгер в контекст
		ctx := WithContext(req.Context(), log)
		c.SetRequest(req.WithContext(ctx))

		// Передаем управление следующему хендлеру
		return next(c)
	}
}

// generateRequestID генерирует уникальный идентификатор запроса (UUID v4)
func GenerateID() string {
	return uuid.New().String()
}
