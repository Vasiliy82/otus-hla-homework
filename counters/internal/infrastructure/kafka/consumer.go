package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/common/utils"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/usecases"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
	useCase  *usecases.ProcessSagaEventUseCase
}

func NewKafkaConsumer(brokers string, groupID string, topic string, useCase *usecases.ProcessSagaEventUseCase) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c, topic: topic, useCase: useCase}, nil
}

func (kc *KafkaConsumer) StartConsuming(ctx context.Context) error {
	if err := kc.consumer.SubscribeTopics([]string{kc.topic}, nil); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			logger.FromContext(ctx).Infow("Consumer closed", "kc.topic", kc.topic)
			return kc.consumer.Close()
		default:
			msg, err := kc.consumer.ReadMessage(100 * time.Millisecond)
			if msg != nil && err == nil {
				newCtx := utils.AddRequestIDToContext(ctx, utils.ExtractRequestIDFromKafka(msg.Headers))
				var sagaEvent domain.SagaEvent
				if err := json.Unmarshal(msg.Value, &sagaEvent); err == nil {
					_ = kc.useCase.Execute(newCtx, sagaEvent)
				}
			}
		}
	}
}

func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}
