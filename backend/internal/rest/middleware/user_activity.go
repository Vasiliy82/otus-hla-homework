package middleware

import (
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/backend/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/apperrors"
	"github.com/labstack/echo/v4"
)

func UserActivityMiddleware(svc domain.SocialNetworkService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем claims из контекста
			claims, ok := c.Get("claims").(*domain.UserClaims)
			if !ok || claims.Subject == "" {
				return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid user claims"))
			}

			// Обновляем last activity
			userID := domain.UserKey(claims.Subject)
			if err := svc.SetLastActivity(userID); err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}

			return next(c)
		}
	}
}
