package tests

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend-messages/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend-messages/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend-messages/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/backend-messages/internal/repository"
	"github.com/Vasiliy82/otus-hla-homework/backend-messages/internal/services"
	"github.com/google/uuid"
)

const (
	configFile = "app-messages-test.yaml"
)

func TestGenerateRandomMessages(t *testing.T) {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Инициализация базы данных
	dbCluster, err := postgresqldb.InitDBCluster(ctx, cfg.SQLServer)
	if err != nil {
		t.Fatalf("Failed to initialize database cluster: %v", err)
	}
	defer dbCluster.Close()

	// Создаем репозиторий и сервис
	dialogRepository := repository.NewDialogRepository(dbCluster)
	dialogService := services.NewDialogService(cfg.Dialogs, dialogRepository)

	// Генерация сообщений
	numMessages := 100000
	numWorkers := 50                  // Указываем количество воркеров
	users := generateRandomUsers(100) // Создаем 100 пользователей

	err = generateRandomMessages(ctx, cancel, dialogService, users, numMessages, numWorkers)
	if err != nil {
		t.Fatalf("Failed to generate random messages: %v", err)
	}

	t.Logf("Successfully generated %d random messages using %d workers", numMessages, numWorkers)
}

// generateRandomUsers генерирует случайных пользователей с валидными UUID
func generateRandomUsers(count int) []domain.UserKey {
	users := make([]domain.UserKey, count)
	for i := 0; i < count; i++ {
		users[i] = domain.UserKey(uuid.NewString())
	}
	return users
}

// generateRandomMessages генерирует случайные сообщения между пользователями с использованием воркеров
func generateRandomMessages(ctx context.Context, cancel context.CancelFunc, dialogService domain.DialogService, users []domain.UserKey, numMessages int, numWorkers int) error {
	// Каналы для распределения задач и получения результатов
	tasks := make(chan int, numMessages)
	results := make(chan error, numMessages)

	// Для синхронизации завершения воркеров
	var wg sync.WaitGroup

	// Запускаем воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, dialogService, users, tasks, results)
	}

	// Генерируем задачи
	go func() {
		for i := 0; i < numMessages; i++ {
			select {
			case <-ctx.Done():
				close(tasks)
				return
			case tasks <- i:
			}
		}
		close(tasks)
	}()

	// Ожидаем завершения работы воркеров
	var firstErr error
	for i := 0; i < numMessages; i++ {
		select {
		case <-ctx.Done():
			// firstErr = ctx.Err()
			break
		case err := <-results:
			if err != nil && firstErr == nil {
				firstErr = err
				cancel() // Отменяем контекст при первой ошибке
				break
			}
		}
	}

	// Ждём завершения всех воркеров
	wg.Wait()
	close(results)

	return firstErr
}

// worker обрабатывает задачи из канала tasks
func worker(ctx context.Context, wg *sync.WaitGroup, dialogService domain.DialogService, users []domain.UserKey, tasks <-chan int, results chan<- error) {
	defer wg.Done()

	for task := range tasks {
		select {
		case <-ctx.Done():
			results <- ctx.Err()
			return
		default:
			// Выбираем случайных пользователей
			sender := users[rand.Intn(len(users))]
			receiver := users[rand.Intn(len(users))]
			// if sender == receiver {
			// 	results <- nil // Пропускаем задачи без ошибок
			// 	continue
			// }

			// Генерируем текст сообщения
			message := fmt.Sprintf("Random message #%d from %s to %s", task+1, sender, receiver)

			// Отправляем сообщение через сервис
			err := dialogService.SendMessage(ctx, sender, receiver, message)
			results <- err
		}
	}
}
