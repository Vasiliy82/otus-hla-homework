package logger

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level string, writer io.Writer, _ ...zap.Option) (*zap.Logger, error) {
	var atomicLevel zap.AtomicLevel
	if level == "" {
		atomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else {
		parsedLevel, errParse := zap.ParseAtomicLevel(level)
		if errParse != nil {
			return nil, fmt.Errorf("zap.ParseAtomicLevel: %writer", errParse)
		}
		atomicLevel = parsedLevel
	}

	cfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	enc := zapcore.NewJSONEncoder(cfg)

	return zap.New(zapcore.NewCore(enc, zapcore.AddSync(writer), atomicLevel)), nil
}
