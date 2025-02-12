package kafka

import (
	"context"
	"encoding/json"

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
			return kc.consumer.Close()
		default:
			msg, err := kc.consumer.ReadMessage(-1)
			if err == nil {
				var sagaEvent domain.SagaEvent
				if err := json.Unmarshal(msg.Value, &sagaEvent); err == nil {
					_ = kc.useCase.Execute(ctx, sagaEvent)
				}
			}
		}
	}
}

func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}
