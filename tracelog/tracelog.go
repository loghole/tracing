package tracelog

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/loghole/tracing/internal/metrics"
)

const (
	traceKey = "trace_id"
	spanKey  = "span_id"
)

type Logger interface {
	TraceID(ctx context.Context) string
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, template string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, template string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, template string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, template string, args ...interface{})
	With(args ...interface{}) Logger
	WithJSON(key string, b []byte) Logger
}

type TraceLogger struct {
	*zap.SugaredLogger
}

func NewTraceLogger(logger *zap.SugaredLogger) *TraceLogger {
	return &TraceLogger{
		SugaredLogger: logger.Desugar().WithOptions(zap.AddCallerSkip(1), zap.Hooks(metricHook)).Sugar(),
	}
}

func (l *TraceLogger) Debug(ctx context.Context, args ...interface{}) {
	l.withSpanContext(ctx).Debug(args...)
}

func (l *TraceLogger) Debugf(ctx context.Context, template string, args ...interface{}) {
	l.withSpanContext(ctx).Debugf(template, args...)
}

func (l *TraceLogger) Info(ctx context.Context, args ...interface{}) {
	l.withSpanContext(ctx).Info(args...)
}

func (l *TraceLogger) Infof(ctx context.Context, template string, args ...interface{}) {
	l.withSpanContext(ctx).Infof(template, args...)
}

func (l *TraceLogger) Warn(ctx context.Context, args ...interface{}) {
	l.withSpanContext(ctx).Warn(args...)
}

func (l *TraceLogger) Warnf(ctx context.Context, template string, args ...interface{}) {
	l.withSpanContext(ctx).Warnf(template, args...)
}

func (l *TraceLogger) Error(ctx context.Context, args ...interface{}) {
	l.withSpanContext(ctx).Error(args...)
	setErrorTag(ctx)
}

func (l *TraceLogger) Errorf(ctx context.Context, template string, args ...interface{}) {
	l.withSpanContext(ctx).Errorf(template, args...)
	setErrorTag(ctx)
}

func (l TraceLogger) With(args ...interface{}) Logger {
	l.SugaredLogger = l.SugaredLogger.With(args...)

	return &l
}

func (l *TraceLogger) WithJSON(key string, b []byte) Logger {
	var obj interface{}

	if err := json.Unmarshal(b, &obj); err != nil {
		return l.With(key, "unmarshal failed", "failed_json", string(b))
	}

	return l.With(key, obj)
}

func (l *TraceLogger) TraceID(ctx context.Context) string {
	return TraceID(ctx)
}

func (l *TraceLogger) withSpanContext(ctx context.Context) *zap.SugaredLogger {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		return l.SugaredLogger.Desugar().With(
			zap.Stringer(traceKey, sc.TraceID()),
			zap.Stringer(spanKey, sc.SpanID()),
		).Sugar()
	}

	return l.SugaredLogger
}

func TraceID(ctx context.Context) string {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		return sc.TraceID().String()
	}

	return ""
}

func setErrorTag(ctx context.Context) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetStatus(codes.Error, "error")
	}
}

func metricHook(entry zapcore.Entry) error { // nolint:gocritic // implement zap.Hooks()
	switch entry.Level { // nolint:exhaustive // used need only this values.
	case zapcore.DebugLevel:
		metrics.DebugLogsCounter.Inc()
	case zapcore.InfoLevel:
		metrics.InfoLogsCounter.Inc()
	case zapcore.WarnLevel:
		metrics.WarnLogsCounter.Inc()
	case zapcore.ErrorLevel:
		metrics.ErrorLogsCounter.Inc()
	}

	return nil
}
