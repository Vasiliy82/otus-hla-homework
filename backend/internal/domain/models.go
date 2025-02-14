package domain

import "context"

type SagaEventType string

// DialogCounter представляет собой счетчик непрочитанных сообщений в диалоге.
type DialogCounter struct {
	DialogID    string `json:"dialog_id"`
	UnreadCount int    `json:"unread_count"`
}

// SagaEvent представляет события для управления транзакциями SAGA.
type SagaEvent struct {
	TransactionID string        `json:"transaction_id"`
	Type          SagaEventType `json:"type"` //
	DialogID      string        `json:"dialog_id"`
	Error         string        `json:"error,omitempty"` // Для ошибок
}

// SagaCoordinator определяет интерфейс взаимодействия между SAGA и DialogService
type SagaCoordinator interface {
	// Фиксация SAGA (из Saga → DialogService)
	// CommitSagaTransaction(ctx context.Context, transactionID string) error

	// Откат SAGA (из Saga → DialogService)
	// RollbackSagaTransaction(ctx context.Context, transactionID string) error

	// Отправка события SAGA (из DialogService → Saga)
	PublishSagaEvent(ctx context.Context, event SagaEvent) error

	HandleSagaEvent(ctx context.Context, event SagaEvent) error
}

const (
	SagaMessageSent            SagaEventType = "MessageSent"
	SagaCounterIncremented     SagaEventType = "CounterIncremented"
	SagaCounterIncrementFailed SagaEventType = "CounterIncrementFailed"
)
