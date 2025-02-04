package services

import (
	"context"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
)

type CacheWarmup struct {
	cfg      *config.CacheConfig
	userRepo domain.UserRepository
	producer *broker.Producer
}

func NewCacheWarmup(cfg *config.CacheConfig, userRepo domain.UserRepository, producer *broker.Producer) *CacheWarmup {
	return &CacheWarmup{
		cfg:      cfg,
		userRepo: userRepo,
		producer: producer,
	}
}

func (c *CacheWarmup) CacheWarmup(ctx context.Context) error {
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
				event := domain.EventFeedChanged{
					UserID: val,
				}
				if err := c.producer.SendFeedChangedEvent(event); err != nil {
					logger.Logger().Warnw("CacheInvalidator.CacheWarmup: Failed to send cache event", "UserId", event.UserID, "error", err)
				}
				logger.Logger().Debugw("CacheInvalidator.CacheWarmup: sent message [domain.EventFeedRefresh]", "UserId", val)
			}
		}
		logger.Logger().Info("CacheInvalidator.CacheWarmup: Cache warmup completed")
	}()

	return nil
}
