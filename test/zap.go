package test

import (
	"bytes"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MockLogger is logger wrapped for tests.
type MockLogger struct {
	*zap.SugaredLogger
	buf bytes.Buffer
}

func NewMockLogger() *MockLogger {
	logger := &MockLogger{}

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     NoopTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	logger.SugaredLogger = zap.New(zapcore.NewCore(encoder, zapcore.AddSync(&logger.buf), zapcore.DebugLevel)).Sugar()

	return logger
}

func (m *MockLogger) String() string {
	return m.buf.String()
}

func NoopTimeEncoder(time.Time, zapcore.PrimitiveArrayEncoder) {}
