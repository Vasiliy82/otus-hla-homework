package globallogger

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"go.uber.org/zap"
)

var (
	globalLogger *zap.SugaredLogger
	once         sync.Once
)

// InitializeGlobalLogger инициализирует глобальный логгер
func InitializeGlobalLogger(level string, writer io.Writer) error {
	var err error
	once.Do(func() {
		zapLogger, initErr := logger.NewLogger(level, writer)
		if initErr != nil {
			err = fmt.Errorf("failed to initialize logger: %w", initErr)
			return
		}
		globalLogger = zapLogger.Sugar()
	})
	return err
}

// Error логирует ошибку уровня error
func Error(ctx context.Context, err error, message string) {
	globalLogger.Errorw(message, "error", err, "context", ctx)
}

// Errorf логирует ошибку уровня error с форматированием
func Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	globalLogger.Errorw(fmt.Sprintf(format, args...), "error", err, "context", ctx)
}

// Error логирует ошибку уровня error
func Warn(ctx context.Context, err error, message string) {
	globalLogger.Warnw(message, "error", err, "context", ctx)
}

// Errorf логирует ошибку уровня error с форматированием
func Warnf(ctx context.Context, err error, format string, args ...interface{}) {
	globalLogger.Warnw(fmt.Sprintf(format, args...), "error", err, "context", ctx)
}

// Info логирует сообщение уровня info
func Info(ctx context.Context, message string) {
	globalLogger.Infow(message, "context", ctx)
}

// Infof логирует сообщение уровня info с форматированием
func Infof(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Infow(fmt.Sprintf(format, args...), "context", ctx)
}

// Debug логирует сообщение уровня debug
func Debug(ctx context.Context, message string) {
	globalLogger.Debugw(message, "context", ctx)
}

// Debugf логирует сообщение уровня debug с форматированием
func Debugf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.Debugw(fmt.Sprintf(format, args...), "context", ctx)
}
