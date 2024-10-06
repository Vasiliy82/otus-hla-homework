package httpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/internal/services/user"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(ctx context.Context, cfg *config.APIConfig, userHandler user.UserHandler, jwtSvc domain.JWTService) error {

	// Start Server
	address := cfg.ServerAddress

	// prepare echo
	e := echo.New()

	// Настройка middleware CORS
	middleware.CORSConfig(e)

	timeout := cfg.ContextTimeout

	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))
	e.Use(counterMiddleware)

	// Роуты
	e.POST("/api/login", userHandler.Login)
	e.POST("/api/user/register", userHandler.RegisterUser)

	// protected routes
	apiGroup := e.Group("/api")
	apiGroup.Use(middleware.JWTMiddleware(jwtSvc))
	apiGroup.GET("/user/get/:id", userHandler.Get)
	apiGroup.GET("/user/search", userHandler.Search)
	apiGroup.POST("/logout", userHandler.Logout)
	apiGroup.PUT("/friend/add/:friend_id", userHandler.AddFriend)
	apiGroup.PUT("/friend/remove/:friend_id", userHandler.RemoveFriend)

	// Добавляем эндпоинт для метрик Prometheus
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	logger.Logger().Infof("Otus HLA Homework server starting at %s", address)

	// Запуск сервера
	err := e.Start(address)
	if err != nil {
		return fmt.Errorf("error while starting server: %w", err)
	}
	return nil
}
