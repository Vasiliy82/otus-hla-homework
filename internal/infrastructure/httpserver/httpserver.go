package httpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/internal/services/user"
	"github.com/labstack/echo/v4"
)

func Start(ctx context.Context, cfg *config.APIConfig, userHandler user.UserHandler) error {

	// Start Server
	address := cfg.ServerAddress

	// prepare echo
	e := echo.New()

	// Настройка middleware CORS
	middleware.CORSConfig(e)

	timeout := cfg.ContextTimeout

	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Роуты
	e.POST("/api/login", userHandler.Login)
	e.POST("/api/user/register", userHandler.RegisterUser)
	e.GET("/api/user/get/:id", userHandler.Get)

	logger.Logger().Infof("Otus HLA Homework server starting at %s", address)

	// Запуск сервера
	err := e.Start(address)
	if err != nil {
		return fmt.Errorf("Error while starting server: %w", err)
	}
	return nil
}
