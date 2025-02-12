package logger

import (
	"context"
	"os"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ключ для хранения логгера в контексте
type contextKey struct{}

var (
	baseLogger *zap.SugaredLogger
	once       sync.Once
)

// InitLogger инициализирует базовый логгер с основными метаданными
func InitLogger(serviceName, instanceID string) *zap.SugaredLogger {
	once.Do(func() {
		config := zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = colorLevelEncoder

		logger := zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			zapcore.AddSync(zapcore.Lock(os.Stdout)),
			zapcore.DebugLevel,
		))

		defer logger.Sync() // flushes buffer, if any

		baseLogger = logger.Sugar().With(
			"service", serviceName,
			"instance_id", instanceID,
		)
	})
	return baseLogger
}

// FromContext извлекает логгер из контекста или возвращает базовый
func FromContext(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(contextKey{}).(*zap.SugaredLogger)
	if !ok {
		return baseLogger
	}
	return logger
}

// WithContext добавляет логгер в контекст
func WithContext(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func Logger() *zap.SugaredLogger {
	return baseLogger
}

// GetFuncName возвращает имя текущей функции
func GetFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	parts := strings.Split(fn.Name(), "/")
	return parts[len(parts)-1]
}

// colorLevelEncoder добавляет цвет в зависимости от уровня логирования
func colorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch level {
	case zapcore.DebugLevel:
		color = "\033[90m" // Gray
	case zapcore.InfoLevel:
		color = "\033[97m" // White
	case zapcore.WarnLevel:
		color = "\033[33m" // Yellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = "\033[97;41m" // White on Dark Red (Brownish) background
	default:
		color = "\033[0m" // Reset
	}

	enc.AppendString(color + level.CapitalString() + "\033[0m")
}
