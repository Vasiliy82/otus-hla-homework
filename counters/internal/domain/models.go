package domain

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

const (
	SagaMessageSent            SagaEventType = "MessageSent"
	SagaCounterIncremented     SagaEventType = "CounterIncremented"
	SagaCounterIncrementFailed SagaEventType = "CounterIncrementFailed"
)
