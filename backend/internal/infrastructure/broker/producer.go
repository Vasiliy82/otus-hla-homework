package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	cfg      *config.Config
}

func NewKafkaProducer(cfg *config.KafkaConfig) (*kafka.Producer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers":  cfg.Brokers,
		"acks":               cfg.Acks,
		"retries":            cfg.Retries,
		"linger.ms":          cfg.LingerMs,
		"enable.idempotence": cfg.EnableIdempotence,
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func NewProducer(ctx context.Context, producer *kafka.Producer, cfg *config.Config) *Producer {
	if producer == nil {
		return nil
	}
	p := &Producer{
		producer: producer,
		cfg:      cfg,
	}
	p.startErrorLogger(ctx)
	return p
}

func (p *Producer) SendPostModifiedEvent(event domain.EventPostModified) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Producer.SendPostModifiedEvent: json.Marshal() returned error: %w", err)
	}

	// Отправляем событие в Kafka
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.cfg.Kafka.TopicPostModified, Partition: kafka.PartitionAny},
		Key:            nil, // round robin distribution
		Value:          eventData,
	}, nil)

	if err != nil {
		return fmt.Errorf("Producer.SendPostModifiedEvent: p.producer.Produce() returned error: %w", err)
	}

	logger.Logger().Debugw("Producer.SendPostModifiedEvent: event was sent", "event", event)

	return nil
}

func (p *Producer) SendFeedChangedEvent(event domain.EventFeedChanged) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Producer.SendFeedChangedEvent: json.Marshal() returned error: %w", err)
	}

	// Отправляем событие в Kafka
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.cfg.Kafka.TopicFeedChanged, Partition: kafka.PartitionAny},
		Key:            []byte(event.UserID),
		Value:          eventData,
	}, nil)

	if err != nil {
		return fmt.Errorf("Producer.SendFeedChangedEvent: p.producer.Produce() returned error: %w", err)
	}

	logger.Logger().Debugw("Producer.SendFeedChangedEvent: event was sent", "event", event)

	return nil
}

func (p *Producer) SendFollowerNotifyEvent(event domain.EventFollowerNotify) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Producer.SendFollowerNotifyEvent: json.Marshal() returned error: %w", err)
	}

	// Отправляем событие в Kafka
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.cfg.Kafka.TopicFollowerNotify, Partition: kafka.PartitionAny},
		Key:            nil, // round robin distribution
		Value:          eventData,
	}, nil)

	if err != nil {
		return fmt.Errorf("Producer.SendFollowerNotifyEvent: p.producer.Produce() returned error: %w", err)
	}

	logger.Logger().Debugw("Producer.SendFollowerNotifyEvent: event was sent", "event", event)

	return nil
}

func (p *Producer) startErrorLogger(ctx context.Context) {
	go func() {
		for {
			select {
			case e := <-p.producer.Events():
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						// Логируем ошибки доставки сообщений
						logger.Logger().Errorw("Producer.StartErrorLogger: p.producer.Events() returned Kafka delivery error",
							"key", string(ev.Key),
							"value", string(ev.Value),
							"error", ev.TopicPartition.Error,
						)
					}
				case kafka.Error:
					// Логируем специфические ошибки Kafka-продюсера
					if ev.IsFatal() {
						logger.Logger().Fatalw("Producer.StartErrorLogger: p.producer.Events() returned Kafka fatal error",
							"code", ev.Code(),
							"message", ev.String(),
						)
					} else {
						logger.Logger().Errorw("Producer.StartErrorLogger: p.producer.Events() returned Kafka error",
							"code", ev.Code(),
							"message", ev.String(),
						)
					}
				}
			case <-ctx.Done():
				logger.Logger().Info("Producer.StartErrorLogger: Stopping Kafka error logger")
				return
			}
		}
	}()
}
