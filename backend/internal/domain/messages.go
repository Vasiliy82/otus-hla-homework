package domain

import (
	"context"
	"net/http"
	"time"
)

type DialogKey string

// DialogMessage представляет сообщение в диалоге
type DialogMessage struct {
	DialogId DialogKey `json:"dialog_id"`
	AuthorId UserKey   `json:"author_id"`
	Datetime time.Time `json:"datetime"`
	Message  string    `json:"message"`
}

// DialogRepository определяет интерфейс репозитория для работы с базой данных
type DialogRepository interface {
	// SaveMessage сохраняет сообщение
	SaveMessage(ctx context.Context, myId, partnerId UserKey, message string) error

	// GetMessages получает сообщения между двумя пользователями
	GetMessages(ctx context.Context, myId, partnerId UserKey, limit, offset int) ([]DialogMessage, error)
}

// DialogService определяет интерфейс для работы с диалогами
type DialogService interface {
	// SendMessage отправляет сообщение
	SendMessage(ctx context.Context, myId, partnerId UserKey, message string) error

	// GetDialog возвращает диалог для пользователя
	GetDialog(ctx context.Context, myID, partnerId UserKey, limit, offset int) ([]DialogMessage, error)
}

// DialogHandler определяет интерфейс обработчика запросов для работы с диалогами
type DialogHandler interface {
	// SendMessage обрабатывает отправку сообщения
	SendMessage(w http.ResponseWriter, r *http.Request)

	// GetDialog обрабатывает получение диалога
	GetDialog(w http.ResponseWriter, r *http.Request)
}
