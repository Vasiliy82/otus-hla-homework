package utils

import (
	"github.com/google/uuid"
)

// generateRequestID генерирует уникальный идентификатор запроса (UUID v4)
func GenerateID() string {
	return uuid.New().String()
}
