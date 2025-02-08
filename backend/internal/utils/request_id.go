package utils

import (
	"context"
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/google/uuid"
)

// GetRequestID извлекает x-request-id из контекста
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(domain.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// AddRequestIDToOutgoing добавляет x-request-id в HTTP header
func AddRequestIDToHeader(ctx context.Context, hdr http.Header) {
	log := logger.FromContext(ctx)
	requestID := GetRequestID(ctx)
	if requestID != "" {

		headerRequestID := hdr.Get(domain.RequestIDHeader)
		if headerRequestID != "" {
			log.Warnw("X-Request-ID already exists", "headerRequestID", headerRequestID)
		}
		hdr.Set("x-request-id", requestID)
	}
}

// generateRequestID генерирует уникальный идентификатор запроса (UUID v4)
func GenerateID() string {
	return uuid.New().String()
}
