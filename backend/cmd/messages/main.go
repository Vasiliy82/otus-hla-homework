package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/tarantool"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/repository"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	"github.com/labstack/echo/v4"
)

const (
	defaultConfigFilename = "socnet.yaml"
	appName               = "messages"
)

func main() {
	log := logger.Logger()

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка системных сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Infof("Received signal: %s, shutting down...", sig)
		cancel()
	}()

	// Читаем путь к файлу конфигурации из флага
	configPath := flag.String("config", defaultConfigFilename, "path to the configuration file")
	flag.Parse()

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Валидируем конфигурацию
	validationErrors := cfg.Validate()
	if len(validationErrors) > 0 {
		fmt.Println("Configuration validation failed with the following errors:")
		for i, vErr := range validationErrors {
			fmt.Printf("%d. %s\n", i+1, vErr.Error())
		}
		log.Fatal("Config validation failure")
	}

	var dialogRepository domain.DialogRepository

	// Инициализация сервисов и обработчиков
	log.Infof("Initializing services...")

	if cfg.Dialogs.UseInmem {
		tarconn, err := tarantool.NewTarConn(*cfg.Tarantool)
		if err != nil {
			log.Fatalf("Error connecting to Tarantool: %v", err)
		}
		dialogRepository = repository.NewdialogRepositoryTar(tarconn)
	} else {
		// Инициализация подключения к базе данных
		log.Infof("Initializing PostgreSQL connection...")

		// Инициализация кластера базы данных
		dbCluster, err := postgresqldb.InitDBCluster(ctx, cfg.SQLServer, appName)
		if err != nil {
			logger.Logger().Fatalw("Failed to initialize database cluster", "err", err)
		}
		defer dbCluster.Close()

		dialogRepository = repository.NewDialogRepository(dbCluster)
	}

	dialogService := services.NewDialogService(cfg.Dialogs, dialogRepository)
	dialogHandler := rest.NewDialogHandler(dialogService)

	// Инициализация Echo
	e := echo.New()
	e.Use(middleware.SetRequestContextWithTimeout(cfg.API.ContextTimeout))
	middleware.CORSConfig(e)

	// Роутинг
	e.POST("/api/dialog/:partnerId/send", dialogHandler.SendMessage)
	e.GET("/api/dialog/:partnerId/list", dialogHandler.GetDialog)
	e.GET("/api/dialog", dialogHandler.GetDialogs)

	// Запуск HTTP-сервера
	go func() {
		log.Infof("Server is running on %s", cfg.API.ServerAddress)
		if err := e.Start(cfg.API.ServerAddress); err != nil && err != context.Canceled {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Ожидаем завершения
	<-ctx.Done()
	log.Infof("Shutting down server...")

	// Завершаем сервер с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.API.ShutdownTimeout)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Infof("Error during server shutdown: %v", err)
	}

}
