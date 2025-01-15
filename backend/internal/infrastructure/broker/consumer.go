package broker

import (
	log "github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// NewKafkaConsumer создает новый экземпляр Kafka Consumer
func NewKafkaConsumer(brokers, groupID, topic string) (*kafka.Consumer, error) {
	log.Logger().Debugw("NewKafkaConsumer()", "brokers", brokers, "groupID", groupID, "topic", topic)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	// Подписка на топик
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
