package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	echo "github.com/labstack/echo/v4"
)

// SetRequestContextWithTimeout устанавливает таймаут для каждого HTTP-запроса и логирует таймауты
func SetRequestContextWithTimeout(d time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			log := logger.FromContext(ctx).With("func", logger.GetFuncName())
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()

			newRequest := c.Request().WithContext(ctx)
			c.SetRequest(newRequest)
			err := next(c)
			select {
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					log.Warnw("Request timed out",
						"path", c.Path(), "method", c.Request().Method,
						"URL", c.Request().URL.String(), "timeout", d)
					return echo.NewHTTPError(http.StatusGatewayTimeout, "Request timed out")
				}
			default:
				// Контекст завершился нормально
			}

			return err
		}
	}
}
