package broker

import (
	log "github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// NewKafkaConsumer создает новый экземпляр Kafka Consumer
func NewKafkaConsumer1(brokers, groupID, topic string) (*kafka.Consumer, error) {
	log.Logger().Debugw("NewKafkaConsumer()", "brokers", brokers, "groupID", groupID, "topic", topic)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     brokers,
		"group.id":              groupID,
		"auto.offset.reset":     "earliest",
		"session.timeout.ms":    6000,  // 6 секунд
		"max.poll.interval.ms":  60000, // 60 секунд
		"heartbeat.interval.ms": 250,   // .25 секунд
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
