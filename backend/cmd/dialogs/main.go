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
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/tarantool"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/repository"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/usecases"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/utils"
	commw "github.com/Vasiliy82/otus-hla-homework/common/infrastructure/http/middleware"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultConfigFilename = "socnet.yaml"
	appName               = "dialogs"
)

func main() {

	// Создание основного контекста с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Инициализация базового логгера
	log := logger.InitLogger(appName, utils.GenerateID())

	// Сохранение логгера в контекст
	ctx = logger.WithContext(ctx, log)
	// Локальный логгер с дополнительным контекстом
	log = log.With("func", logger.GetFuncName())

	log.Debug("Started")

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
			log.Fatalw("Failed to initialize database cluster", "err", err)
		}
		defer dbCluster.Close()

		dialogRepository = repository.NewDialogRepository(dbCluster)
	}

	kafkaProd, err := broker.NewKafkaProducer(cfg.Kafka)

	prod := broker.NewProducer(ctx, kafkaProd, cfg)
	if err != nil {
		log.Fatalf("Error init producer: %v", err)
	}

	dialogService := services.NewDialogService(cfg.Dialogs, dialogRepository)
	dialogHandler := rest.NewDialogHandler(dialogService)

	sagaOrch := usecases.NewSagaOrchestrator(prod, dialogService)

	dialogService.SetSagaCoordinator(sagaOrch)

	sagaProc := services.NewSagaBusProcessor(cfg, sagaOrch)
	sagaProc.Start(ctx)

	metrics := commw.NewPrometheusMetrics()
	metrics.Register() // Регистрируем метрики в Prometheus

	// Инициализация Echo
	e := echo.New()
	e.Use(commw.RequestIDMiddleware)

	// Добавляем middleware для сбора метрик RED
	e.Use(commw.PrometheusMetricsMiddleware(metrics))

	// Отдельный эндпоинт для метрик Prometheus
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.Use(middleware.SetRequestContextWithTimeout(cfg.API.ContextTimeout))

	// Роутинг
	e.GET("/api/dialog", dialogHandler.GetDialogs)
	e.GET("/api/dialog/:partnerId/list", dialogHandler.GetDialog)
	e.POST("/api/dialog/:partnerId/send", dialogHandler.SendMessage)

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
	log.Infof("Waiting Kafka Saga processor")
	sagaProc.Wait()
}
