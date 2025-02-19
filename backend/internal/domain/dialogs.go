package domain

import (
	"context"
	"net/http"
	"time"
)

type DialogKey string
type MessageKey int

type Dialog struct {
	UserId   UserKey   `json:"user_id"`
	DialogId DialogKey `json:"dialog_id"`
}

// DialogMessage представляет сообщение в диалоге
type DialogMessage struct {
	DialogId  DialogKey  `json:"dialog_id"`
	MessageId MessageKey `json:"message_id"`
	AuthorId  UserKey    `json:"author_id"`
	Datetime  time.Time  `json:"datetime"`
	Message   string     `json:"message"`
}

// DialogRepository определяет интерфейс репозитория для работы с базой данных
type DialogRepository interface {
	// SaveMessage сохраняет сообщение
	// SaveMessage(ctx context.Context, myId, partnerId UserKey, message string) error

	// SaveMessageWithSaga сохраняет сообщение с TransactionID и SagaStatus
	SaveMessageWithSaga(ctx context.Context, myId, partnerId UserKey, message, transactionID string) error

	// GetDialog получает сообщения между двумя пользователями
	GetDialog(ctx context.Context, myId, partnerId UserKey, limit, offset int) ([]DialogMessage, error)

	// GetDialogs получает список диалогов для данного пользователя
	GetDialogs(ctx context.Context, myId UserKey, limit, offset int) ([]Dialog, error)

	// UpdateSagaStatus обновляет статус SAGA (pending → committed/failed)
	UpdateSagaStatus(ctx context.Context, transactionID string, status string) error
}

// DialogService определяет интерфейс для работы с диалогами и SAGA
type DialogService interface {
	// SendMessage отправляет сообщение и возвращает TransactionID (используется в SAGA)
	SendMessage(ctx context.Context, myId, partnerId UserKey, message string) (string, error)

	// CommitSagaTransaction фиксирует транзакцию SAGA (подтверждает сообщение)
	CommitSagaTransaction(ctx context.Context, transactionID string) error

	// RollbackSagaTransaction откатывает транзакцию SAGA (удаляет сообщение)
	RollbackSagaTransaction(ctx context.Context, transactionID string) error

	// GetDialog возвращает диалог для пользователя
	GetDialog(ctx context.Context, myID, partnerId UserKey, limit, offset int) ([]DialogMessage, error)

	// GetDialogs получает список диалогов для данного пользователя
	GetDialogs(ctx context.Context, myId UserKey, limit, offset int) ([]Dialog, error)

	SetSagaCoordinator(sagaCoordinator SagaCoordinator)
}

// DialogHandler определяет интерфейс обработчика запросов для работы с диалогами
type DialogHandler interface {
	// SendMessage обрабатывает отправку сообщения
	SendMessage(w http.ResponseWriter, r *http.Request)

	// GetDialog обрабатывает получение диалога
	GetDialog(w http.ResponseWriter, r *http.Request)

	// GetDialog получает список диалогов для данного пользователя
	GetDialogs(w http.ResponseWriter, r *http.Request)
}
