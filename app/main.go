package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	log "github.com/Vasiliy82/otus-hla-homework/internal/observability/globallogger"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	pgRepo "github.com/Vasiliy82/otus-hla-homework/internal/repository/postgres"

	"github.com/Vasiliy82/otus-hla-homework/internal/rest"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/internal/service"
	"github.com/joho/godotenv"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Errorf(context.Background(), err, "Error loading .env file")
	}
}

func main() {
	// Инициализация глобального логгера
	err := log.InitializeGlobalLogger("debug", os.Stdout)
	if err != nil {
		panic("failed to initialize global logger")
	}

	// Prepare database
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	db, err := sql.Open(`postgres`, dsn)
	if err != nil {
		log.Errorf(context.Background(), err, "failed to open connection to database")
		return
	}
	err = db.Ping()
	if err != nil {
		log.Errorf(context.Background(), err, "failed to ping database")
		return
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Errorf(context.Background(), err, "got error when closing the DB connection")
		}
	}()

	// Инициализация сервисов
	userService := service.NewUserService(pgRepo.NewUserRepository(db), pgRepo.NewSessionRepository(db))
	userHandler := rest.NewUserHandler(userService)

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	// prepare echo
	e := echo.New()
	e.Use(middleware.CORS)
	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Warnf(context.Background(), err, "failed to parse timeout, using default timeout! %v", err)
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Роуты
	e.POST("/login", userHandler.Login)
	e.POST("/user/register", userHandler.RegisterUser)
	e.GET("/user/get/:id", userHandler.Get)

	log.Infof(context.Background(), "Otus HLA Homework server starting at %s", address)

	// Запуск сервера
	err = e.Start(address)
	if err != nil {
		log.Errorf(context.Background(), err, "Error while starting server")
	}
}
