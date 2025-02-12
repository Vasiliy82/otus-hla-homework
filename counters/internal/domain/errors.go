// internal/domain/errors.go
package domain

import "errors"

var (
	ErrDialogNotFound    = errors.New("dialog not found")
	ErrCounterOverflow   = errors.New("counter overflow")
	ErrProcessingMessage = errors.New("error processing message")
)

// // ErrorEvent представляет собой структуру для ошибок, отправляемых в Kafka.
// type ErrorEvent struct {
// 	DialogID string `json:"dialog_id"`
// 	Reason   string `json:"reason"`
// }
