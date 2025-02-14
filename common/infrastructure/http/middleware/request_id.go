package middleware

import (
	"github.com/Vasiliy82/otus-hla-homework/common/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/utils"
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
			requestID = utils.GenerateID()
		}

		// Сохранение x-request-id в контекст
		ctx = utils.AddRequestIDToContext(ctx, requestID)
		// ctx = logger.WithContext(context.WithValue(ctx, domain.contextKey{}, requestID), logger.FromContext(ctx).With("requestID", requestID))

		c.SetRequest(req.WithContext(ctx))

		// Добавление x-request-id в заголовок ответа
		utils.AddRequestIDToHeader(ctx, c.Response().Header())

		err := next(c)

		return err

	}
}
