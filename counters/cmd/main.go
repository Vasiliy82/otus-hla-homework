package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal"
	config "github.com/Vasiliy82/otus-hla-homework/counters/internal/config-counters"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/utils"
)

const (
	defaultConfigFilename = "counters.yaml"
	appName               = "counters"
)

func main() {

	log := logger.InitLogger(appName, utils.GenerateID())

	// Флаг для указания пути к файлу конфигурации
	configPath := flag.String("config", defaultConfigFilename, "path to the configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Валидируем конфигурацию
	validationErrors := cfg.Validate()
	if len(validationErrors) > 0 {
		log.Error("Configuration validation failed with the following errors:")
		for i, vErr := range validationErrors {
			log.Errorf("%d. %s\n", i+1, vErr.Error())
		}
		log.Fatal("Config validation failure")
	}

	log.Info("Configuration validation successful!")

	// Создаем контекст с обработкой сигналов завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Инициализация приложения (все зависимости и конфигурации)
	app, err := internal.InitializeApp(ctx, cfg, appName)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Запуск приложения (HTTP сервер и Kafka Consumer)
	if err := app.Run(ctx, cfg); err != nil {
		log.Fatalf("application stopped with error: %v", err)
	}

	log.Info("application shutdown gracefully")
}
