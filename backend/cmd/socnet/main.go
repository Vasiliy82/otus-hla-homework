package main

import (
	"context"
	"flag"
	"fmt"
	"net/http/httputil"
	"net/url"
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
	appName               = "socnet"
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

	// Валидируем конфигурацию
	validationErrors := cfg.Validate()
	if len(validationErrors) > 0 {
		fmt.Println("Configuration validation failed with the following errors:")
		for i, vErr := range validationErrors {
			fmt.Printf("%d. %s\n", i+1, vErr.Error())
		}
		log.Logger().Fatal("Config validation failure")
	}

	// log.Logger().Debugw("Config", "cfg", cfg)

	log.Logger().Debug("main: init postgresql...")
	db, err := postgresqldb.InitDBCluster(ctx, cfg.SQLServer, appName)
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
	prod, err := broker.NewKafkaProducer(cfg.Kafka)
	if err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
	bprod := broker.NewProducer(ctx, prod, cfg)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init User Service...")
	snService := services.NewSocialNetworkService(cfg, userRepo, postRepo, postCache, jwtService, bprod)
	log.Logger().Debugln("done")
	log.Logger().Debug("main: init User Handler...")
	userHandler := rest.NewSocialNetworkHandler(snService, cfg.API)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init Event processor (Post modified) ...")
	procPostModified := services.NewPostModifiedProcessor(cfg, snService, bprod)
	procPostModified.Start(ctx)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init Event processor (Feed Changed) ...")
	procFeedChanged := services.NewFeedChangedProcessor(cfg, postRepo, postCache)
	procFeedChanged.Start(ctx)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: Warm up cache...")
	cacheWarmup := services.NewCacheWarmup(cfg.Cache, userRepo, bprod)
	cacheWarmup.CacheWarmup(ctx)
	log.Logger().Debugln("done")

	log.Logger().Debug("main: init reverse proxy (messages)...")
	rpMessagesTarget, err := url.Parse(cfg.SocialNetwork.SvcMessagesURL)
	if err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
	rpMessages := httputil.NewSingleHostReverseProxy(rpMessagesTarget)
	log.Logger().Debugln("done")

	log.Logger().Debugln("Starting HTTP server...")
	err = httpserver.Start(ctx, cfg, userHandler, jwtService, snService, rpMessages, jwtService)
	cancel()
	procFeedChanged.Wait()
	procPostModified.Wait()

	if err != nil {
		log.Logger().Fatalf("Error: %v", err)
	}
	log.Logger().Debug("main: main: otus-hla-homework succesfully stopped")
}
