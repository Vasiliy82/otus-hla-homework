package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type PostModifiedProcessor struct {
	w         broker.Worker
	cfg       *config.Config
	snService domain.SocialNetworkService
	producer  *broker.Producer
}

func NewPostModifiedProcessor(cfg *config.Config, snService domain.SocialNetworkService, producer *broker.Producer) *PostModifiedProcessor {
	return &PostModifiedProcessor{
		cfg:       cfg,
		snService: snService,
		producer:  producer,
	}
}

func (p *PostModifiedProcessor) Start(ctx context.Context) {
	p.w = *broker.NewWorker(&broker.WorkerConfig{
		Topic:         p.cfg.Kafka.TopicPostModified,
		NumWorkers:    p.cfg.Kafka.NumWorkersPostModified,
		FuncProcessor: p.process,
		ConsumerConfig: &kafka.ConfigMap{
			"bootstrap.servers":     p.cfg.Kafka.Brokers,
			"enable.auto.commit":    false,
			"group.id":              p.cfg.Kafka.CGPostModified,
			"auto.offset.reset":     "earliest",
			"session.timeout.ms":    6000,  // 6 секунд
			"max.poll.interval.ms":  60000, // 60 секунд
			"heartbeat.interval.ms": 2000,  // 1 секунд
		},
	})
	p.w.Start(ctx)
}

func (p *PostModifiedProcessor) process(ctx context.Context, msg *kafka.Message, workerId int) error {
	var event domain.EventPostModified
	var err error

	if err = json.Unmarshal(msg.Value, &event); err != nil {
		logger.Logger().Warnw("PostModifiedProcessor.process: json.Unmarshal() returned error", "workerId", workerId, "err", err)
		logger.Logger().Debugw("PostModifiedProcessor.process: json.Unmarshal() returned error", "workerId", workerId, "msg.Value", msg.Value)
		return err
	}
	logger.Logger().Debugw("Event processing started", "event", event)
	ids := event.AntiLadyGagaIds
	if len(ids) == 0 {
		ids, err = p.snService.GetFriendsIds(event.Post.UserId)
		if err != nil {
			return fmt.Errorf("PostModifiedProcessor.process:  p.snService.GetFriendsIds() returned error: %w", err)
		}
	} else {
		logger.Logger().Debugw("Lady gaga post-processing", "ids", ids)
	}
	if len(ids) > p.cfg.SocialNetwork.MaxPostCreatedPerWorker {
		logger.Logger().Debugw("Lady gaga detected", "len(ids)", len(ids))
		for len(ids) > p.cfg.SocialNetwork.PostCreatedPacketSize {
			chunkSize := p.cfg.SocialNetwork.PostCreatedPacketSize
			if chunkSize > len(ids) {
				chunkSize = len(ids)
			}
			chunk := ids[:chunkSize]
			ids = ids[chunkSize:]
			pm := domain.EventPostModified{
				Event:           event.Event,
				Post:            event.Post,
				AntiLadyGagaIds: chunk,
			}
			if err = p.producer.SendPostModifiedEvent(pm); err != nil {
				logger.Logger().Warnf("PostModifiedProcessor.process: p.producer.SendPostModifiedEvent() returned error: %v", err)
			}
		}
		return nil
	}

	for _, id := range ids {
		fn := domain.EventFollowerNotify{
			Recipient: id,
			Content: &domain.EventFollowerNotifyContent{
				Post:  event.Post,
				Event: event.Event,
			},
		}
		fc := domain.EventFeedChanged{
			UserID: id,
		}
		p.producer.SendFollowerNotifyEvent(fn)
		p.producer.SendFeedChangedEvent(fc)
	}

	return nil
}

func (p *PostModifiedProcessor) Wait() {
	p.w.Wait()
}
