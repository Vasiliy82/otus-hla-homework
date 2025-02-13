package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	commw "github.com/Vasiliy82/otus-hla-homework/common/infrastructure/http/middleware"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	httpAdapter "github.com/Vasiliy82/otus-hla-homework/counters/internal/adapters/http"
	config "github.com/Vasiliy82/otus-hla-homework/counters/internal/config-counters"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/infrastructure/kafka"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/infrastructure/repository"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// App структура для хранения зависимостей и управления жизненным циклом приложения
type App struct {
	httpServer    *echo.Echo
	kafkaConsumer *kafka.KafkaConsumer
	pool          *pgxpool.Pool
}

// InitializeApp конфигурирует зависимости и возвращает готовое приложение
func InitializeApp(ctx context.Context, cfg *config.ConfigCounters, appName string) (*App, error) {
	// Инициализация PostgreSQL
	pool, err := initDBInstance(ctx, cfg.Postgres, appName)
	if err != nil {
		return nil, err
	}
	pgRepo := repository.NewPGCounterRepository(pool)

	// Инициализация Redis
	redisRepo := repository.NewRedisCounterRepository(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.Db)

	// Создание комбинированного репозитория для использования Redis и PostgreSQL
	combinedRepo := repository.NewCombinedCounterRepository(redisRepo, pgRepo)

	// Инициализация Kafka Producer
	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.TopicSagaBus)
	if err != nil {
		return nil, err
	}

	// Создание Services, Use cases
	counterService := usecases.NewCounterService(combinedRepo)
	processEventUC := usecases.NewProcessSagaEventUseCase(counterService, kafkaProducer)

	// Настройка HTTP сервера
	httpServer := initializeHTTPServer(counterService)

	// Инициализация Kafka Consumer
	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.CGSagaBus, cfg.Kafka.TopicSagaBus, processEventUC)
	if err != nil {
		return nil, err
	}

	return &App{
		httpServer:    httpServer,
		kafkaConsumer: kafkaConsumer,
		pool:          pool,
	}, nil
}

// Run запускает HTTP сервер и Kafka Consumer
func (a *App) Run(ctx context.Context, cfg *config.ConfigCounters) error {
	// Запуск HTTP сервера
	go func() {
		if err := a.httpServer.Start(cfg.API.ServerAddress); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start HTTP server on address %s: %v", cfg.API.ServerAddress, err)
		}
	}()

	// Запуск Kafka Consumer
	go func() {
		if err := a.kafkaConsumer.StartConsuming(ctx); err != nil {
			log.Fatalf("failed to consume Kafka messages: %v", err)
		}
	}()

	// Ожидание завершения
	<-ctx.Done()

	// Graceful Shutdown
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	a.pool.Close()

	a.kafkaConsumer.Close()

	return nil
}

func initDBInstance(ctx context.Context, cfg *config.DBInstanceConfig, appName string) (*pgxpool.Pool, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	val := url.Values{}
	val.Add("sslmode", "disable")
	val.Add("application_name", appName)
	dsn := fmt.Sprintf("%s?%s", connStr, val.Encode())

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MinConns = int32(cfg.MinConns)
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	log.Infof("connection pool established: %s", connStr)

	return pool, nil
}

// initializeHTTPServer настройка HTTP роутов и серверов
func initializeHTTPServer(counterService domain.CounterService) *echo.Echo {
	e := echo.New()
	e.Use(commw.RequestIDMiddleware)
	handler := httpAdapter.NewCounterHTTPHandler(counterService)
	handler.RegisterRoutes(e)
	return e
}
