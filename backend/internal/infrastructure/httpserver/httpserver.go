package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/wsreverseproxy"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	cmw "github.com/Vasiliy82/otus-hla-homework/common/infrastructure/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(ctx context.Context,
	cfg *config.Config,
	snHandler services.SocialNetworkHandler,
	jwtSvc domain.JWTService,
	snSvc domain.SocialNetworkService,
	jwtService domain.JWTService,
) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())

	// Start Server
	address := cfg.API.ServerAddress

	// prepare router
	router, err := rest.NewProxyRouter(cfg.SocialNetwork)
	if err != nil {
		return fmt.Errorf("httpserver.Start: rest.NewProxyRouter returned error: %w", err)
	}

	// prepare echo
	e := echo.New()

	e.Use(cmw.RequestIDMiddleware)

	// Настройка middleware CORS
	middleware.CORSConfig(e)

	timeout := cfg.API.ContextTimeout

	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	counterMiddleware := NewMetricsMiddleware(prometheus.DefaultRegisterer, *cfg.Metrics)

	e.Use(counterMiddleware)

	// Публикация статического файла index.html
	e.File("/", "./frontend-demo/index.html")

	// Роуты
	e.POST("/api/login", snHandler.Login)
	e.POST("/api/user/register", snHandler.CreateUser)

	e.GET("/ws", wsreverseproxy.WebSocketHandler(cfg.SocialNetwork.SvcPostsWsURL, jwtService))

	// protected routes
	apiGroup := e.Group("/api")
	apiGroup.Use(middleware.JWTMiddleware(jwtSvc))
	apiGroup.Use(middleware.UserActivityMiddleware(snSvc))
	apiGroup.GET("/user/get/:id", snHandler.GetUser)
	apiGroup.GET("/user/search", snHandler.Search)
	apiGroup.POST("/logout", snHandler.Logout)
	apiGroup.PUT("/friend/add/:friend_id", snHandler.AddFriend)
	apiGroup.PUT("/friend/remove/:friend_id", snHandler.RemoveFriend)

	apiGroup.POST("/post", snHandler.CreatePost)
	apiGroup.GET("/post", snHandler.ListPosts)
	apiGroup.GET("/post/:post_id", snHandler.GetPost)
	apiGroup.PUT("/post/:post_id", snHandler.UpdatePost)
	apiGroup.DELETE("/post/:post_id", snHandler.DeletePost)
	apiGroup.GET("/post/feed", snHandler.GetFeed)
	routerHandler := router.RouterHandler(cfg.SocialNetwork)
	// Проксирование маршрутов для микросервиса сообщений
	apiGroup.Any("/*", routerHandler)

	// Добавляем эндпоинт для метрик Prometheus
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// проксирование прочих маршрутов
	e.Any("/*", routerHandler, middleware.JWTMiddleware(jwtSvc))

	log.Infof("Otus HLA Homework server starting at %s", address)

	// Запуск сервера в горутине
	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			log.Errorf("error while starting server: %v", err)
		}
	}()

	// Ожидание завершения контекста
	<-ctx.Done()
	log.Info("Shutting down server...")

	// Установка таймаута для завершения
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.API.ShutdownTimeout)*time.Second)
	defer cancel()

	// Корректное завершение сервера с использованием таймаута
	if err := e.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error while shutting down server: %w", err)
	}

	log.Info("Server shutdown gracefully")
	return nil
}
