package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

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
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//prepare database
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
		log.Fatal("failed to open connection to database", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

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
		log.Println("failed to parse timeout, using default timeout!")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	e.POST("/login", userHandler.Login)
	e.POST("/user/register", userHandler.RegisterUser)
	e.GET("/user/get/:id", userHandler.Get)

	log.Println("Otus HLA Homework")

	e.Logger.Fatal(e.Start(address))
}
