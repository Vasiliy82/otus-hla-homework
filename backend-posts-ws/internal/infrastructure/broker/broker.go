package broker

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

// NewKafkaConsumer создает Kafka Consumer с поддержкой кастомного логирования через zap
func NewKafkaConsumer(brokers, groupID, topic string, logger *zap.SugaredLogger) (*kafka.Consumer, error) {
	// Создаем Consumer с дополнительными настройками логирования
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":      brokers,
		"group.id":               groupID,
		"auto.offset.reset":      "earliest",
		"go.logs.channel.enable": true, // Включаем передачу логов в канал
		"log_level":              3,
	})
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
		return nil, err
	}

	// Запускаем обработку логов в отдельной горутине
	go handleKafkaLogs(consumer, logger)

	// Подписываемся на топик
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logger.Fatalf("Failed to subscribe to topics: %v", err)
		return nil, err
	}

	return consumer, nil
}

// handleKafkaLogs читает логи Kafka из канала и записывает их через zap
func handleKafkaLogs(consumer *kafka.Consumer, logger *zap.SugaredLogger) {
	for logEvent := range consumer.Logs() {
		// Записываем логи в формате JSON
		logger.Info("KafkaLog",
			zap.String("name", logEvent.Name),
			zap.String("tag", logEvent.Tag),
			zap.String("message", logEvent.Message),
			zap.Int("level", logEvent.Level),
		)
	}
}
