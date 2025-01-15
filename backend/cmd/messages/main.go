package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/repository"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest/middleware"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	"github.com/labstack/echo/v4"
)

const (
	defaultConfigFilename = "socnet.yaml"
)

func main() {

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка системных сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s, shutting down...", sig)
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

	// Инициализация подключения к базе данных
	log.Println("Initializing database connection...")

	// Инициализация кластера базы данных
	dbCluster, err := postgresqldb.InitDBCluster(ctx, cfg.SQLServer)
	if err != nil {
		logger.Logger().Fatalw("Failed to initialize database cluster", "err", err)
	}
	defer dbCluster.Close()

	// Использование пула базы данных
	if err := timeQuery(ctx, dbCluster); err != nil {
		logger.Logger().Errorw("Error in exampleQuery", "err", err)
	}

	// Инициализация сервисов и обработчиков
	log.Println("Initializing services...")
	dialogRepository := repository.NewDialogRepository(dbCluster)
	dialogService := services.NewDialogService(cfg.Dialogs, dialogRepository)
	dialogHandler := rest.NewDialogHandler(dialogService)

	// Инициализация Echo
	e := echo.New()
	e.Use(middleware.SetRequestContextWithTimeout(cfg.API.ContextTimeout))
	middleware.CORSConfig(e)

	// Роутинг
	e.POST("/dialog/:partnerId/send", dialogHandler.SendMessage)
	e.GET("/dialog/:partnerId/list", dialogHandler.GetDialog)

	// Запуск HTTP-сервера
	go func() {
		log.Printf("Server is running on %s", cfg.API.ServerAddress)
		if err := e.Start(cfg.API.ServerAddress); err != nil && err != context.Canceled {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Ожидаем завершения
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Завершаем сервер с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.API.ShutdownTimeout)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

}

func timeQuery(ctx context.Context, cluster *postgresqldb.DBCluster) error {
	// Получаем пул для чтения
	dbPool, err := cluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return err
	}

	// Выполняем запрос на получение времени из базы данных
	query := "SELECT NOW()"
	var dbTime time.Time
	err = dbPool.QueryRow(ctx, query).Scan(&dbTime)
	if err != nil {
		return err
	}

	// Локальное время
	localTime := time.Now()

	// Рассчитываем разницу
	timeDiff := localTime.Sub(dbTime)
	logMessage := "Local time and database server time are in sync."

	// Проверяем, если разница превышает 1 секунду, фиксируем это как проблему
	if timeDiff > time.Second || timeDiff < -time.Second {
		logMessage = "WARNING: Local time and database server time differ!"
		logger.Logger().Warnw(logMessage, "local_time", localTime, "db_time", dbTime, "time_diff", timeDiff)
	} else {
		logger.Logger().Infow(logMessage, "local_time", localTime, "db_time", dbTime, "time_diff", timeDiff)
	}

	log.Printf("%s Local: %s, DB: %s, Diff: %s", logMessage, localTime, dbTime, timeDiff)
	return nil
}
