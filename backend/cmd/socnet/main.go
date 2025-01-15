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
	log "github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/repository/cache"
	"github.com/redis/go-redis/v9"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/repository"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/httpserver"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
)

const (
	defaultTimeout        = 30
	defaultAddress        = ":9090"
	defaultConfigFilename = "socnet.yaml"
)

func main() {
	var jwtService domain.JWTService
	var err error

	log.Logger().Debug("main: starting otus-hla-homework")

	// Создание основного контекста с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Создаем канал для получения системных сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Горутинa для обработки системных сигналов
	go func() {
		sig := <-sigs
		log.Logger().Infof("Received signal: %s, shutting down...\n", sig)
		cancel() // Отмена контекста при получении сигнала
	}()

	// Флаг для указания пути к файлу конфигурации
	configPath := flag.String("config", defaultConfigFilename, "path to the configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Logger().Fatalf("Error loading config: %v", err)
	}

	// log.Logger().Debugw("Config", "cfg", cfg)

	log.Logger().Debug("main: init postgresql...")
	db, err := postgresqldb.InitDBCluster(ctx, cfg.SQLServer)
	if err != nil {
		log.Logger().Errorf("main: postgresqldb.InitDB returned error: %v", err)
		return
	}

	defer db.Close()

	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port), // Адрес Redis сервера (по умолчанию localhost:6379)
		Password: cfg.Cache.Redis.Password,                                         // Пароль, если установлен
		DB:       0,                                                                // Номер базы данных (по умолчанию 0)
	})

	// Проверяем подключение к Redis
	_, err = redis.Ping(context.Background()).Result()
	if err != nil {
		log.Logger().Fatalw("Ошибка подключения к Redis:", "err", err)
		return
	}

	log.Logger().Debugln("done")

	log.Logger().Debug("main: init metrics...")
	postgresqldb.StartMonitoring(db, cfg.Metrics.UpdateInterval)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init JWT Service...")
	if jwtService, err = services.NewJWTService(cfg.JWT, repository.NewBlacklistRepository(ctx, db)); err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init User Repository...")
	userRepo := repository.NewUserRepository(ctx, db)
	log.Logger().Debugln("done")
	log.Logger().Debug("main: init Post Repository...")
	postRepo := repository.NewPostRepository(ctx, db)
	log.Logger().Debugln("done")
	log.Logger().Debug("main: init Post Cache...")
	postCache := cache.NewPostCache(cfg.Cache, redis)
	log.Logger().Debugln("done")
	log.Logger().Debug("main: init Producer...")
	prod, err := broker.NewKafkaProducer(cfg.Cache.Kafka)
	if err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
	bprod := broker.NewProducer(prod, cfg.Cache.InvalidateTopic)
	bprod.StartErrorLogger(ctx)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init Cache Invalidator...")
	cacheInv := services.NewCacheInvalidator(cfg.Cache, cfg.SocialNetwork, userRepo, postRepo, postCache, bprod)
	cacheInv.ListenAndProcessEvents(ctx)
	cacheInv.CacheWarmup(ctx)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init User Service...")
	snService := services.NewSocialNetworkService(cfg, userRepo, postRepo, postCache, jwtService, bprod)
	log.Logger().Debugln("done")
	log.Logger().Debug("main: init User Handler...")
	userHandler := rest.NewSocialNetworkHandler(snService, cfg.API)
	log.Logger().Debugln("done")

	log.Logger().Debugln("Starting HTTP server...")
	err = httpserver.Start(ctx, cfg, userHandler, jwtService, snService)
	if err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
	cacheInv.WaitForDone()
	log.Logger().Debug("main: main: otus-hla-homework succesfully stopped")
}
