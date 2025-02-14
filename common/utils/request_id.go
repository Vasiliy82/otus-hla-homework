package utils

import (
	"context"
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/common/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

// ContextKey для хранения x-request-id в контексте
type contextKey struct{}

// generateRequestID генерирует уникальный идентификатор запроса (UUID v4)
func GenerateID() string {
	return uuid.New().String()
}

// ExtractRequestIDFromContext извлекает X-Request-ID из контекста
func ExtractRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(contextKey{}).(string); ok {
		return requestID
	}
	return "" // Если нет, генерируем новый
}

// AddRequestIDToContext добавляет X-Request-ID в контекст
func AddRequestIDToContext(ctx context.Context, requestID string) context.Context {

	oldRequestId := ExtractRequestIDFromContext(ctx)
	ctx = context.WithValue(ctx, contextKey{}, requestID)

	if oldRequestId != requestID {
		newLogger := logger.FromContext(ctx).With(domain.RequestIDVariable, requestID)
		ctx = logger.WithContext(ctx, newLogger)
	}

	return ctx
}

// ExtractRequestIDFromKafka извлекает X-Request-ID из заголовков Kafka
func ExtractRequestIDFromKafka(headers []kafka.Header) string {
	for _, h := range headers {
		if h.Key == domain.RequestIDHeader {
			return string(h.Value)
		}
	}
	return "" // Если нет, генерируем новый
}

// AddRequestIDToKafka добавляет X-Request-ID в заголовки Kafka
func AddRequestIDToKafka(ctx context.Context, headers *[]kafka.Header) {
	requestID := ExtractRequestIDFromContext(ctx)
	if requestID != "" {
		*headers = append(*headers, kafka.Header{Key: domain.RequestIDHeader, Value: []byte(requestID)})
	}
}

// // AddRequestIDToHeader добавляет X-Request-ID в HTTP заголовки ответа
// func AddRequestIDToEchoHeader(ctx context.Context, header echo.Header) {
// 	requestID := ExtractRequestIDFromContext(ctx)
// 	header.Set(domain.RequestIDHeader, requestID)
// }

// AddRequestIDToOutgoing добавляет x-request-id в HTTP header
func AddRequestIDToHeader(ctx context.Context, hdr http.Header) {
	requestID := ExtractRequestIDFromContext(ctx)
	if requestID != "" {
		hdr.Set(domain.RequestIDHeader, requestID)
	}

}
