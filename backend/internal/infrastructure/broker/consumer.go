package broker

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// NewKafkaConsumer создает новый экземпляр Kafka Consumer
func NewKafkaConsumer(brokers, groupID, topic string) (*kafka.Consumer, error) {
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
