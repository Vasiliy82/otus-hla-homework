package middleware

import (
	"net/http"
	"strings"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtService domain.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Logger().Debug("middleware.JWTMiddleware: missing or invalid Authorization header", "err")
				return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid Authorization header"))
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Валидация токена
			token, err := jwtService.ValidateToken(domain.TokenString(tokenString))
			if err != nil {
				logger.Logger().Warnw("middleware.JWTMiddleware: token validation error", "err", err, "tokenString", tokenString)
				return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("invalid token"))
			}

			claims, err := jwtService.ExtractClaims(token)
			if err != nil {
				logger.Logger().Warnw("middleware.JWTMiddleware: jwtService.ExtractClaims() returned error", "err", err, "token", token)
				return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("invalid token"))
			}

			logger.Logger().Debugw("Token successfully validated", "user_id", claims.Subject)

			// Сохранение токена в контексте
			c.Set("token", token)
			c.Set("claims", claims)

			return next(c)
		}
	}
}
