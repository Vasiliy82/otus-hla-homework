package middleware

import (
	"net/http"
	"net/http/httputil"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/labstack/echo/v4"
)

// ProxyMiddleware создает middleware для проксирования запросов через ReverseProxy
func ReverseProxyMiddleware(proxy *httputil.ReverseProxy) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлечение claims из контекста
			claims, ok := c.Get("claims").(*domain.UserClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid claims in context")
			}

			// Извлечение user_id (Subject) из claims
			userID := claims.Subject
			if userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing user ID in claims")
			}

			// Добавляем заголовок с ID пользователя
			c.Request().Header.Set("X-User-Id", userID)

			// Используем прокси для обработки запроса
			proxy.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	}
}
