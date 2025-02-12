// internal/domain/services.go
package domain

import "context"

// CounterService определяет интерфейс для работы с счетчиками.
type CounterService interface {
	IncrementUnread(ctx context.Context, dialogID string) error
	ResetUnread(ctx context.Context, dialogID string) error
}

// CounterRepository определяет интерфейс для взаимодействия с хранилищем счетчиков.
type CounterRepository interface {
	GetCounter(ctx context.Context, dialogID string) (*DialogCounter, error)
	SaveCounter(ctx context.Context, counter *DialogCounter) error
}

// KafkaPublisher определяет интерфейс для публикации событий в Kafka.
type KafkaPublisher interface {
	PublishSagaEvent(ctx context.Context, event SagaEvent) error
}

// KafkaConsumer определяет интерфейс для подписки на события Kafka.
type KafkaConsumer interface {
	StartConsuming(ctx context.Context) error
}
