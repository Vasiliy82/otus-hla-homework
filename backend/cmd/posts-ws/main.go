package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap/zapcore"
)

const (
	defaultConfigFilename = "socnet.yaml"
)

func main() {
	// Инициализация логгера
	logger.SetLogger(logger.NewStdOut(nil))
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

	logger.SetLevel(getLogLevel(cfg.Log.Level))

	// Настраиваем Kafka Consumer
	// consumer, err := broker.NewKafkaConsumer(cfg.Broker.Brokers, cfg.Broker.Group, cfg.Broker.Topic, log)
	// if err != nil {
	// 		log.Fatalf("Failed to create Kafka consumer: %v", err)
	// }
	// defer consumer.Close()

	m := middleware.NewPrometheusMetrics()
	m.Register()

	// Настраиваем Echo
	e := echo.New()

	e.Use(middleware.ZapLoggerMiddleware(log))

	e.Use(middleware.PrometheusMetricsMiddleware(m))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Публикация статического файла index.html
	e.File("/", "./frontend-demo/index.html")

	// Инициализация WebSocket сервера
	wsService := services.NewWSService(cfg.Posts)
	e.GET("/ws", wsService.HandleConnection)

	// Запускаем HTTP-сервер
	go func() {
		log.Infof("Starting WebSocket server on %s...", cfg.API.ServerAddress)
		if err := e.Start(cfg.API.ServerAddress); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	fcProc := services.NewFollowerNotifyProcessor(cfg, wsService)
	fcProc.Start(ctx)

	// Ожидаем завершения
	<-ctx.Done()
	log.Info("Shutting down server...")

	// Завершаем сервер с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.API.ShutdownTimeout)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Errorf("Error during server shutdown: %v", err)
	}
	log.Info("Waiting consumer threads...")
	fcProc.Wait()
}

func getLogLevel(l log.Lvl) zapcore.Level {

	// Настройка уровня логирования
	logLevel := zapcore.DebugLevel
	switch l {
	case log.DEBUG:
		logLevel = zapcore.DebugLevel
	case log.INFO:
		logLevel = zapcore.InfoLevel
	case log.WARN:
		logLevel = zapcore.WarnLevel
	case log.ERROR:
		logLevel = zapcore.ErrorLevel
	case log.OFF:
		logLevel = zapcore.FatalLevel
	}

	return logLevel
}
