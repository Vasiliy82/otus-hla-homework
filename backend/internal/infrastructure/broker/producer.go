package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
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

func NewProducer(producer *kafka.Producer, topic string) *Producer {
	if producer == nil {
		return nil
	}
	return &Producer{
		producer: producer,
		topic:    topic,
	}
}

func (p *Producer) SendCacheEvent(event domain.EventInvalidateCache) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Producer.SendCacheEvent: json.Marshal() returned error: %w", err)
	}

	// Отправляем событие в Kafka
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            []byte(event.UserID),
		Value:          eventData,
	}, nil)

	if err != nil {
		return fmt.Errorf("Producer.SendCacheEvent: p.producer.Produce() returned error: %w", err)
	}

	logger.Logger().Debugw("Producer.SendCacheEvent: event was sent", "event", event)

	return nil
}

func (p *Producer) StartErrorLogger(ctx context.Context) {
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
