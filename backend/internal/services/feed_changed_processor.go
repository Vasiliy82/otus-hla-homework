package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type FeedChangedProcessor struct {
	w         *broker.WorkerPool
	cfg       *config.Config
	postRepo  domain.PostRepository
	postCache domain.PostCache
}

func NewFeedChangedProcessor(cfg *config.Config, postRepo domain.PostRepository,
	postCache domain.PostCache,
) *FeedChangedProcessor {
	return &FeedChangedProcessor{
		cfg:       cfg,
		postRepo:  postRepo,
		postCache: postCache,
	}
}
func (p *FeedChangedProcessor) Start(ctx context.Context) {
	p.w = broker.NewWorker(
		broker.NewWorkerConfig(p.cfg.Kafka.TopicFeedChanged,
			p.cfg.Kafka.NumWorkersFeedChanged,
			p.process,
			&kafka.ConfigMap{
				"bootstrap.servers":     p.cfg.Kafka.Brokers,
				"enable.auto.commit":    false,
				"group.id":              p.cfg.Kafka.CGFeedChanged,
				"auto.offset.reset":     "earliest",
				"session.timeout.ms":    6000,  // 6 секунд
				"max.poll.interval.ms":  60000, // 60 секунд
				"heartbeat.interval.ms": 2000,  // 1 секунд
			}),
	)
	p.w.Start(ctx)
}

func (p *FeedChangedProcessor) process(ctx context.Context, msg *kafka.Message, workerId int) error {
	var event domain.EventFeedChanged
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.Logger().Warnw("FeedChangedProcessor.process: json.Unmarshal() returned error", "workerId", workerId, "err", err)
		logger.Logger().Debugw("FeedChangedProcessor.process: json.Unmarshal() returned error", "workerId", workerId, "msg.Value", msg.Value)
		return err
	}
	logger.Logger().Debugw("Event processing started", "event", event)
	if err := p.recalculateUserFeed(ctx, event.UserID); err != nil {
		logger.Logger().Warnw("FeedChangedProcessor.process: c.recalculateUserFeed() returned error", "userId", event.UserID, "workerId", workerId, "err", err)
		return err
	}
	return nil
}

// recalculateUserFeed пересчитывает ленту указанного пользователя
func (p *FeedChangedProcessor) recalculateUserFeed(ctx context.Context, userId domain.UserKey) error {
	_ = ctx // FIXME:
	posts, err := p.postRepo.GetFeed(userId, p.cfg.SocialNetwork.FeedLength)
	if err != nil {
		logger.Logger().Warnw("CacheInvalidator.recalculateUserFeed: c.postRepo.GetFeed returned error", "userId", userId, "err", err)
		return fmt.Errorf("CacheInvalidator.recalculateUserFeed: c.postRepo.GetFeed returned error: %w", err)
	}
	if p.postCache != nil {
		if err = p.postCache.UpdateFeed(userId, posts); err != nil {
			logger.Logger().Warnw("CacheInvalidator.recalculateUserFeed: c.postCache.UpdateFeed returned error", "userId", userId, "err", err)
			return fmt.Errorf("CacheInvalidator.recalculateUserFeed: c.postCache.UpdateFeed returned error: %w", err)
		}
	}
	logger.Logger().Debugw("CacheInvalidator.recalculateUserFeed: лента пользователя пересчитана", "userId", userId)
	return nil
}

func (p *FeedChangedProcessor) Wait() {
	p.w.Wait()
}
