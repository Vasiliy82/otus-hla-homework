package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	log "github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	_ "github.com/lib/pq"

	"github.com/Vasiliy82/otus-hla-homework/internal/repository"

	"github.com/Vasiliy82/otus-hla-homework/internal/infrastructure/httpserver"
	"github.com/Vasiliy82/otus-hla-homework/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest"
	"github.com/Vasiliy82/otus-hla-homework/internal/services/jwt"
	service "github.com/Vasiliy82/otus-hla-homework/internal/services/user"
	"github.com/joho/godotenv"
)

const (
	defaultTimeout        = 30
	defaultAddress        = ":9090"
	defaultConfigFilename = "app.yaml"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Logger().Errorf("Error loading .env file: %v", err)
	}
}

func main() {
	var jwtService domain.JWTService
	var err error

	// Создание основного контекста с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Создаем канал для получения системных сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Горутинa для обработки системных сигналов
	go func() {
		sig := <-sigs
		fmt.Printf("Received signal: %s, shutting down...\n", sig)
		cancel() // Отмена контекста при получении сигнала
	}()

	// Флаг для указания пути к файлу конфигурации
	configPath := flag.String("config", defaultConfigFilename, "path to the configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Logger().Fatalf("Error loading config: %v", err)
	}

	db, err := postgresqldb.InitDB(ctx, cfg.SQLServer)
	if err != nil {
		log.Logger().Errorf("main: postgresqldb.InitDB returned error: %v", err)
		return
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Logger().Errorf("main: db.Close() returned error: %v", err)
		}
	}()

	// Инициализация сервисов
	if jwtService, err = jwt.NewJWTService(cfg.JWT, repository.NewBlacklistRepository(db)); err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	userService := service.NewUserService(userRepo, sessionRepo, jwtService)
	userHandler := rest.NewUserHandler(userService)

	_ = jwtService

	err = httpserver.Start(ctx, cfg.API, userHandler)
	if err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
}
