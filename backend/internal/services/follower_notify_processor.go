package services

import (
	"context"
	"encoding/json"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type FollowerNotifyProcessor struct {
	w   broker.WorkerPool
	cfg *config.Config
	ws  *WSService
}

func NewFollowerNotifyProcessor(cfg *config.Config, wsService *WSService) *FollowerNotifyProcessor {
	return &FollowerNotifyProcessor{
		cfg: cfg,
		ws:  wsService,
	}
}
func (p *FollowerNotifyProcessor) Start(ctx context.Context) {
	w := broker.NewWorker(
		broker.NewWorkerConfig(
			p.cfg.Kafka.TopicFollowerNotify,
			p.cfg.Kafka.NumWorkersFollowerNotify,
			p.process,
			&kafka.ConfigMap{
				"bootstrap.servers":     p.cfg.Kafka.Brokers,
				"enable.auto.commit":    false,
				"group.id":              p.cfg.Kafka.CGFollowerNotify,
				"auto.offset.reset":     "earliest",
				"session.timeout.ms":    6000,  // 6 секунд
				"max.poll.interval.ms":  60000, // 60 секунд
				"heartbeat.interval.ms": 2000,  // 1 секунд
			},
		),
	)
	w.Start(ctx)
}

func (p *FollowerNotifyProcessor) process(ctx context.Context, msg *kafka.Message, workerId int) error {
	var event domain.EventFollowerNotify
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.Logger().Warnw("FollowerNotifyProcessor.process: json.Unmarshal() returned error", "workerId", workerId, "err", err)
		logger.Logger().Debugw("FollowerNotifyProcessor.process: json.Unmarshal() returned error", "workerId", workerId, "msg.Value", msg.Value)
		return err
	}
	logger.Logger().Debugw("FollowerNotifyProcessor.process: event processing started", "event", event)
	if err := p.sendToWebsocket(event); err != nil {
		logger.Logger().Warnw("FollowerNotifyProcessor.process: p.sendToWebsocket() returned error", "userId", event.Recipient, "workerId", workerId, "err", err)
		return err
	}
	return nil
}

func (p *FollowerNotifyProcessor) sendToWebsocket(event domain.EventFollowerNotify) error {
	return p.ws.SendFollowerNotifyEvent(event.Recipient, event.Content)
}

func (p *FollowerNotifyProcessor) Wait() {
	p.w.Wait()
}
