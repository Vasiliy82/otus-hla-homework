package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ZapLoggerMiddleware возвращает middleware для логирования запросов с использованием zap
func ZapLoggerMiddleware(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Выполняем обработку запроса
			err := next(c)

			// Логируем детали запроса
			logger.Infow("HTTP request",
				"method", c.Request().Method,
				"path", c.Path(),
				"status", c.Response().Status,
				"latency", time.Since(start).String(),
				"remote_ip", c.RealIP(),
				"user_agent", c.Request().UserAgent(),
			)

			return err
		}
	}
}
