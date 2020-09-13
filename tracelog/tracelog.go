package tracelog

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/gadavy/tracing/internal"
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
		SugaredLogger: logger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
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
	setErrorTag(ctx)
	l.withSpanContext(ctx).Error(args...)
}

func (l *TraceLogger) Errorf(ctx context.Context, template string, args ...interface{}) {
	setErrorTag(ctx)
	l.withSpanContext(ctx).Errorf(template, args...)
}

func (l TraceLogger) With(args ...interface{}) Logger {
	l.SugaredLogger = l.SugaredLogger.With(args...)

	return &l
}

func (l *TraceLogger) WithJSON(key string, b []byte) Logger {
	var obj interface{}

	if err := jsoniter.Unmarshal(b, &obj); err != nil {
		return l.With(key, "unmarshal failed", "failed_json", string(b))
	}

	return l.With(key, obj)
}

func (l *TraceLogger) TraceID(ctx context.Context) string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if sc, ok := span.Context().(internal.SpanContext); ok {
			return sc.TraceID().String()
		}
	}

	return ""
}

func (l *TraceLogger) withSpanContext(ctx context.Context) *zap.SugaredLogger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if sc, ok := span.Context().(internal.SpanContext); ok {
			return l.Desugar().With(
				zap.Stringer(traceKey, sc.TraceID()),
				zap.Stringer(spanKey, sc.SpanID()),
			).Sugar()
		}
	}

	return l.SugaredLogger
}

func setErrorTag(ctx context.Context) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ext.Error.Set(span, true)
	}
}
