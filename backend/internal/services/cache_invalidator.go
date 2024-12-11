package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// CacheInvalidator сервис для пересчета кешей лент новостей
type CacheInvalidator struct {
	cfg       *config.CacheConfig
	snCfg     *config.SocialNetworkConfig
	postCache domain.PostCache
	postRepo  domain.PostRepository
	userRepo  domain.UserRepository
	producer  *broker.Producer
	wg        sync.WaitGroup
}

// NewCacheInvalidator создает новый экземпляр CacheInvalidator
func NewCacheInvalidator(cfg *config.CacheConfig, snCfg *config.SocialNetworkConfig, userRepo domain.UserRepository, postRepo domain.PostRepository, postCache domain.PostCache, producer *broker.Producer) *CacheInvalidator {
	return &CacheInvalidator{
		cfg:       cfg,
		snCfg:     snCfg,
		postCache: postCache,
		postRepo:  postRepo,
		userRepo:  userRepo,
		producer:  producer,
	}
}

// ListenAndProcessEvents запускает numWorkers параллельных воркеров для обработки событий из Kafka.
func (c *CacheInvalidator) ListenAndProcessEvents(ctx context.Context) error {
	// Запуск numWorkers горутин с отдельными Consumer для каждого
	for i := 0; i < c.cfg.InvalidateNumWorkers; i++ {
		c.wg.Add(1)
		workerId := i
		go func() {
			defer c.wg.Done()
			if err := c.startWorker(ctx, workerId); err != nil {
				logger.Logger().Errorw("CacheInvalidator.ListenAndProcessEvents: Ошибка в работе воркера", "workerId", workerId, "err", err)
			}
		}()
	}

	return nil
}

// startWorker создает отдельный Kafka Consumer и выполняет обработку сообщений в рамках одного воркера
func (c *CacheInvalidator) startWorker(ctx context.Context, workerId int) error {
	// Создаем Kafka Consumer для данного воркера
	consumer, err := broker.NewKafkaConsumer(c.cfg.Kafka.Brokers, c.cfg.InvalidateConsumerGroup, c.cfg.InvalidateTopic)
	if err != nil {
		return fmt.Errorf("CacheInvalidator.startWorker: %w", err)
	}
	defer consumer.Close()

	logger.Logger().Infow("CacheInvalidator.startWorker: Воркер запущен и подключен к Kafka", "workerId", workerId)

	for {
		select {
		case <-ctx.Done():
			// Завершаем работу воркера, если контекст отменен
			logger.Logger().Infow("CacheInvalidator.startWorker: Контекст отменен, завершаем работу воркера", "workerId", workerId)
			return nil
		default:
			// Извлекаем сообщение из Kafka
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Обрабатываем сообщение
				if err := c.processMessage(ctx, e, workerId); err != nil {
					logger.Logger().Errorw("CacheInvalidator.startWorker: c.processMessage() returned error", "workerId", workerId, "message", e, "err", err)
					// В случае ошибки - повторяем цикл , не подтверждая сообщение
					continue
				}
				// В случае успеха - подтверждаем сообщение в Kafka
				consumer.CommitMessage(e)
			case kafka.Error:
				logger.Logger().Errorw("CacheInvalidator.startWorker: consumer.Poll(100) returned Kafka error", "workerId", workerId, "err", e)
			default:
				// Игнорируем другие события Kafka
				logger.Logger().Debugw("CacheInvalidator.startWorker: consumer.Poll(100) returned unknown Kafka message", "workerId", workerId, "event", e)
			}
		}
	}
}

// processMessage выполняет бизнес-логику обработки сообщения
func (c *CacheInvalidator) processMessage(ctx context.Context, msg *kafka.Message, workerId int) error {
	var event domain.EventInvalidateCache
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.Logger().Warnw("CacheInvalidator.processMessage: json.Unmarshal() returned error", "workerId", workerId, "err", err)
		logger.Logger().Debugw("CacheInvalidator.processMessage: json.Unmarshal() returned error", "workerId", workerId, "msg.Value", msg.Value)
		return err
	}

	switch event.EventType {
	case domain.EventPostCreated, domain.EventPostEdited, domain.EventPostDeleted:
		// Генерируем события для пересчета ленты всех друзей пользователя
		if err := c.generateFeedRefreshEvents(ctx, event.UserID); err != nil {
			logger.Logger().Debugw("CacheInvalidator.processMessage: c.generateFeedRefreshEvents() returned error", "userId", event.UserID, "workerId", workerId, "err", err)
			return err
		}
	case domain.EventFeedRefresh:
		// Выполняем пересчет ленты для указанного пользователя
		if err := c.recalculateUserFeed(ctx, event.UserID); err != nil {
			logger.Logger().Debugw("CacheInvalidator.processMessage: c.recalculateUserFeed() returned error", "userId", event.UserID, "workerId", workerId, "err", err)
			return err
		}
	default:
		logger.Logger().Warnw("CacheInvalidator.processMessage: received unknown event type", "eventType", event.EventType)
		logger.Logger().Debugw("CacheInvalidator.processMessage: received unknown event type", "event", event)
		return fmt.Errorf("CacheInvalidator.processMessage: unknown event: %s", event.EventType)
	}

	return nil
}

// generateFeedRefreshEvents отправляет события пересчета ленты для друзей указанного пользователя
func (c *CacheInvalidator) generateFeedRefreshEvents(ctx context.Context, userId domain.UserKey) error {
	_ = ctx // FIXME:
	// Получаем список друзей пользователя из userRepo
	friends, err := c.userRepo.GetFriendsIds(userId)
	if err != nil {
		return err
	}

	// Генерируем событие для каждого друга
	for _, friendId := range friends {
		if err := c.sendInvalidationEvent(userId, domain.EventFeedRefresh); err != nil {
			logger.Logger().Warnw("CacheInvalidator.generateFeedRefreshEvents: c.sendInvalidationEvent() returned error", "friendId", friendId, "err", err)
		}
	}

	logger.Logger().Debugw("CacheInvalidator.generateFeedRefreshEvents: разосланы события domain.EventFeedRefresh", "userId", userId, "friends", friends)

	return nil
}

// recalculateUserFeed пересчитывает ленту указанного пользователя
func (c *CacheInvalidator) recalculateUserFeed(ctx context.Context, userId domain.UserKey) error {
	_ = ctx // FIXME:
	posts, err := c.postRepo.GetFeed(userId, c.snCfg.FeedLength)
	if err != nil {
		logger.Logger().Warnw("CacheInvalidator.recalculateUserFeed: c.postRepo.GetFeed returned error", "userId", userId, "err", err)
		return fmt.Errorf("CacheInvalidator.recalculateUserFeed: c.postRepo.GetFeed returned error: %w", err)
	}
	if c.postCache != nil {
		if err = c.postCache.UpdateFeed(userId, posts); err != nil {
			logger.Logger().Warnw("CacheInvalidator.recalculateUserFeed: c.postCache.UpdateFeed returned error", "userId", userId, "err", err)
			return fmt.Errorf("CacheInvalidator.recalculateUserFeed: c.postCache.UpdateFeed returned error: %w", err)
		}
	}
	logger.Logger().Debugw("CacheInvalidator.recalculateUserFeed: лента пользователя пересчитана", "userId", userId)
	return nil
}

func (c *CacheInvalidator) WaitForDone() {
	// Ожидание завершения всех воркеров
	c.wg.Wait()
}

func (c *CacheInvalidator) sendInvalidationEvent(userId domain.UserKey, eventType domain.EventType) error {
	if c.producer == nil {
		return nil
	}

	// Создаем и сериализуем событие
	event := domain.EventInvalidateCache{
		UserID:    userId,
		EventType: eventType,
	}

	// Отправляем событие в Kafka
	if err := c.producer.SendCacheEvent(event); err != nil {
		logger.Logger().Errorw("CacheInvalidator.sendInvalidationEvent: s.producer.SendEvent() returned error", "userId", userId, "err", err)
		return err
	}

	return nil
}

func (c *CacheInvalidator) CacheWarmup(ctx context.Context) error {
	if !c.cfg.CacheWarmupEnabled {
		return nil
	}

	// Запускаем прогрев кеша в отдельной горутине
	go func() {

		logger.Logger().Infow("CacheInvalidator.CacheWarmup: Starting cache warmup")

		// Выполнение прогрева
		ids, err := c.userRepo.GetUsersActiveSince(c.cfg.CacheWarmupPeriod)
		if err != nil {
			logger.Logger().Errorw("CacheInvalidator.CacheWarmup: Failed to fetch active users", "error", err)
			return
		}

		for _, val := range ids {
			select {
			case <-ctx.Done():
				logger.Logger().Info("CacheInvalidator.CacheWarmup: Stopping cache warmup due to context cancellation")
				return
			default:
				event := domain.EventInvalidateCache{
					UserID:    val,
					EventType: domain.EventFeedRefresh,
				}
				if err := c.producer.SendCacheEvent(event); err != nil {
					logger.Logger().Warnw("CacheInvalidator.CacheWarmup: Failed to send cache event", "UserId", event.UserID, "error", err)
				}
				logger.Logger().Debugw("CacheInvalidator.CacheWarmup: sent message [domain.EventFeedRefresh]", "UserId", val)
			}
		}
		logger.Logger().Info("CacheInvalidator.CacheWarmup: Cache warmup completed")
	}()

	return nil
}
