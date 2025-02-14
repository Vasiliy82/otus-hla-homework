package services

import (
	"context"
	"encoding/json"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type SagaBusProcessor struct {
	w         *broker.WorkerPool
	cfg       *config.Config
	sagaCoord domain.SagaCoordinator
}

func NewSagaBusProcessor(cfg *config.Config, sagaCoord domain.SagaCoordinator,
) *SagaBusProcessor {
	return &SagaBusProcessor{
		cfg:       cfg,
		sagaCoord: sagaCoord,
	}
}
func (p *SagaBusProcessor) Start(ctx context.Context) {
	p.w = broker.NewWorker(
		broker.NewWorkerConfig(p.cfg.Kafka.TopicSagaBus,
			p.cfg.Kafka.NumWorkersSagaBus,
			p.process,
			&kafka.ConfigMap{
				"bootstrap.servers":     p.cfg.Kafka.Brokers,
				"enable.auto.commit":    false,
				"group.id":              p.cfg.Kafka.CGSagaBus,
				"auto.offset.reset":     "earliest",
				"session.timeout.ms":    6000,  // 6 секунд
				"max.poll.interval.ms":  60000, // 60 секунд
				"heartbeat.interval.ms": 2000,  // 1 секунд
			}),
	)
	p.w.Start(ctx)
}

func (p *SagaBusProcessor) process(ctx context.Context, msg *kafka.Message, workerId int) error {
	var event domain.SagaEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.Logger().Warnw("json.Unmarshal() returned error", "workerId", workerId, "err", err)
		logger.Logger().Debugw("json.Unmarshal() returned error", "workerId", workerId, "msg.Value", msg.Value)
		return err
	}
	logger.Logger().Debugw("Event processing started", "event", event)
	if err := p.sagaCoord.HandleSagaEvent(ctx, event); err != nil {
		logger.Logger().Warnw("c.recalculateUserFeed() returned error", "event", event, "workerId", workerId, "err", err)
		return err
	}
	return nil
}

func (p *SagaBusProcessor) Wait() {
	p.w.Wait()
}
