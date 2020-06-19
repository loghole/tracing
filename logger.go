package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type key string

const (
	_ctxActionKey key    = "action"
	ActionKey     string = "action"
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
	logger *zap.SugaredLogger
}

func NewTraceLogger(logger *zap.SugaredLogger) *TraceLogger {
	return &TraceLogger{logger: logger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()}
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

func (l *TraceLogger) withAction(ctx context.Context) *zap.SugaredLogger {
	if val, ok := ctx.Value(_ctxActionKey).(string); ok {
		return l.logger.With(ActionKey, val)
	}

	return l.logger
}

func withErrorTag(ctx context.Context) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("error", true)
	}
}
