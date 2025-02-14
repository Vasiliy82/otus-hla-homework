package kafka

import (
	"context"
	"encoding/json"

	"github.com/Vasiliy82/otus-hla-homework/common/utils"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(brokers string, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: p, topic: topic}, nil
}

func (kp *KafkaProducer) PublishSagaEvent(ctx context.Context, event domain.SagaEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}
	var headers []kafka.Header
	utils.AddRequestIDToKafka(ctx, &headers)

	return kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          message,
		Headers:        headers,
	}, nil)
}

func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
