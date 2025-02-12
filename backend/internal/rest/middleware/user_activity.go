package middleware

import (
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/labstack/echo/v4"
)

func UserActivityMiddleware(svc domain.SocialNetworkService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем claims из контекста
			claims, ok := c.Get("claims").(*domain.UserClaims)
			if !ok || claims.Subject == "" {
				logger.Logger().Errorw("rest.UserActivityMiddleware: missing or invalid user claims")
				return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid user claims"))
			}

			// Обновляем last activity
			userID := domain.UserKey(claims.Subject)
			if err := svc.SetLastActivity(userID); err != nil {
				logger.Logger().Errorw("rest.UserActivityMiddleware: svc.SetLastActivity() returned error", "err", err)
				return c.JSON(http.StatusInternalServerError, err)
			}

			return next(c)
		}
	}
}
