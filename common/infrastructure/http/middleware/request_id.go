package middleware

import (
	"context"

	"github.com/Vasiliy82/otus-hla-homework/common/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/common/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequestIDMiddleware извлекает x-request-id из заголовков и сохраняет в контекст
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()

		// Извлечение x-request-id из заголовка или генерация нового
		requestID := req.Header.Get(domain.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Сохранение x-request-id в контекст
		ctx = logger.WithContext(context.WithValue(ctx, domain.RequestIDKey, requestID), logger.FromContext(ctx).With("requestID", requestID))

		c.SetRequest(req.WithContext(ctx))

		// Добавление x-request-id в заголовок ответа
		utils.AddRequestIDToHeader(ctx, c.Response().Header())

		err := next(c)

		return err

	}
}
