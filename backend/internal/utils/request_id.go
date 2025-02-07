package utils

import (
	"context"
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/rest/middleware"
	"github.com/google/uuid"
)

// GetRequestID извлекает x-request-id из контекста
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// AddRequestIDToOutgoing добавляет x-request-id в исходящий HTTP-запрос
func AddRequestIDToOutgoing(ctx context.Context, req *http.Request) {
	requestID := GetRequestID(ctx)
	if requestID != "" {
		req.Header.Set("x-request-id", requestID)
	}
}

// generateRequestID генерирует уникальный идентификатор запроса (UUID v4)
func GenerateID() string {
	return uuid.New().String()
}
