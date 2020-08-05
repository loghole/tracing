package tracing

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type key string

const (
	ActionKey         = "trace_id"
	_ctxActionKey key = "action"
)

func ContextWithAction(ctx context.Context, actionID string) context.Context {
	return context.WithValue(ctx, _ctxActionKey, actionID)
}

func ActionFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(_ctxActionKey).(string); ok {
		return val
	}

	return ""
}

type TraceLogger struct {
	actionKey string
	*zap.SugaredLogger
}

func NewTraceLogger(actionKey string, logger *zap.SugaredLogger) *TraceLogger {
	return &TraceLogger{
		SugaredLogger: logger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
		actionKey:     actionKey,
	}
}

func (l *TraceLogger) Debug(ctx context.Context, args ...interface{}) {
	l.withAction(ctx).Debug(args...)
}

func (l *TraceLogger) Debugf(ctx context.Context, template string, args ...interface{}) {
	l.withAction(ctx).Debugf(template, args...)
}

func (l *TraceLogger) Info(ctx context.Context, args ...interface{}) {
	l.withAction(ctx).Info(args...)
}

func (l *TraceLogger) Infof(ctx context.Context, template string, args ...interface{}) {
	l.withAction(ctx).Infof(template, args...)
}

func (l *TraceLogger) Warn(ctx context.Context, args ...interface{}) {
	l.withAction(ctx).Warn(args...)
}

func (l *TraceLogger) Warnf(ctx context.Context, template string, args ...interface{}) {
	l.withAction(ctx).Warnf(template, args...)
}

func (l *TraceLogger) Error(ctx context.Context, args ...interface{}) {
	withErrorTag(ctx)
	l.withAction(ctx).Error(args...)
}

func (l *TraceLogger) Errorf(ctx context.Context, template string, args ...interface{}) {
	withErrorTag(ctx)
	l.withAction(ctx).Errorf(template, args...)
}

func (l TraceLogger) With(args ...interface{}) *TraceLogger {
	l.SugaredLogger = l.SugaredLogger.With(args...)

	return &l
}

func (l *TraceLogger) WithJSON(key string, b []byte) *TraceLogger {
	var obj interface{}
	if err := jsoniter.Unmarshal(b, &obj); err != nil {
		return l.With(
			key, "unmarshal failed",
			"failed_json", string(b),
		)
	}

	return l.With(key, obj)
}

func (l *TraceLogger) withAction(ctx context.Context) *zap.SugaredLogger {
	if val, ok := ctx.Value(_ctxActionKey).(string); ok {
		return l.SugaredLogger.With(l.actionKey, val)
	}

	return l.SugaredLogger
}

func withErrorTag(ctx context.Context) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("error", true)
	}
}
