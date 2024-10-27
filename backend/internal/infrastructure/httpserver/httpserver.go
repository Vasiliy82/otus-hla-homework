package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(ctx context.Context, cfg *config.Config, userHandler services.UserHandler, postHandler services.PostHandler, jwtSvc domain.JWTService) error {

	// Start Server
	address := cfg.API.ServerAddress

	// prepare echo
	e := echo.New()

	// Настройка middleware CORS
	middleware.CORSConfig(e)

	timeout := cfg.API.ContextTimeout

	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	counterMiddleware := NewMetricsMiddleware(prometheus.DefaultRegisterer, *cfg.Metrics)

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

	apiGroup.POST("/post", postHandler.Create)
	apiGroup.GET("/post", postHandler.List)
	apiGroup.GET("/post/:post_id", postHandler.Get)
	apiGroup.PUT("/post/:post_id", postHandler.Update)
	apiGroup.DELETE("/post/:post_id", postHandler.Delete)
	apiGroup.GET("/post/feed", postHandler.Feed)

	// Добавляем эндпоинт для метрик Prometheus
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	logger.Logger().Infof("Otus HLA Homework server starting at %s", address)

	// Запуск сервера в горутине
	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			logger.Logger().Errorf("error while starting server: %v", err)
		}
	}()

	// Ожидание завершения контекста
	<-ctx.Done()
	logger.Logger().Info("Shutting down server...")

	// Установка таймаута для завершения
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.API.ShutdownTimeout)*time.Second)
	defer cancel()

	// Корректное завершение сервера с использованием таймаута
	if err := e.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error while shutting down server: %w", err)
	}

	logger.Logger().Info("Server shutdown gracefully")
	return nil
}
