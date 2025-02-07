package middleware

import (
	"context"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ContextKey для хранения x-request-id в контексте
type contextKey string

const RequestIDKey contextKey = "x-request-id"

// RequestIDMiddleware извлекает x-request-id из заголовков и сохраняет в контекст
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()

		// Извлечение x-request-id из заголовка или генерация нового
		requestID := req.Header.Get("x-request-id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Сохранение x-request-id в контекст
		ctx = logger.WithContext(context.WithValue(ctx, RequestIDKey, requestID), logger.FromContext(ctx).With("requestID", requestID))

		c.SetRequest(req.WithContext(ctx))

		// // Добавление x-request-id в заголовок ответа для трассировки
		// c.Response().Header().Set("x-request-id", requestID)

		err := next(c)

		if c.Response().Header().Get("x-request-id") == "" {
			c.Response().Header().Set("x-request-id", requestID)
		}

		return err

	}
}
