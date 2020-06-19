package tracing

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
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

type BaseLogger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

type TraceLogger struct {
	logger BaseLogger
}

func NewTraceLogger(logger BaseLogger) *TraceLogger {
	return &TraceLogger{logger: logger}
}

func (l *TraceLogger) Debug(ctx context.Context, args ...interface{}) {
	l.logger.Debug(argsWithAction(ctx, args)...)
}

func (l *TraceLogger) Debugf(ctx context.Context, template string, args ...interface{}) {
	l.logger.Debugf(templateWithAction(ctx, template), args...)
}

func (l *TraceLogger) Info(ctx context.Context, args ...interface{}) {
	l.logger.Info(argsWithAction(ctx, args)...)
}

func (l *TraceLogger) Infof(ctx context.Context, template string, args ...interface{}) {
	l.logger.Infof(templateWithAction(ctx, template), args...)
}

func (l *TraceLogger) Warn(ctx context.Context, args ...interface{}) {
	l.logger.Warn(argsWithAction(ctx, args)...)
}

func (l *TraceLogger) Warnf(ctx context.Context, template string, args ...interface{}) {
	l.logger.Warnf(templateWithAction(ctx, template), args...)
}

func (l *TraceLogger) Error(ctx context.Context, args ...interface{}) {
	withErrorTag(ctx)
	l.logger.Error(argsWithAction(ctx, args)...)
}

func (l *TraceLogger) Errorf(ctx context.Context, template string, args ...interface{}) {
	withErrorTag(ctx)
	l.logger.Errorf(templateWithAction(ctx, template), args...)
}

func argsWithAction(ctx context.Context, args []interface{}) []interface{} {
	if action := withAction(ctx); action != "" {
		// Подставляем первым аргументом action ID
		args = append([]interface{}{action}, args...)
	}

	return args
}

func templateWithAction(ctx context.Context, template string) string {
	if action := withAction(ctx); action != "" {
		// Добавляем перед template action ID
		return strings.Join([]string{action, template}, "")
	}

	return template
}

func withAction(ctx context.Context) string {
	if val, ok := ctx.Value(_ctxActionKey).(string); ok {
		return strings.Join([]string{"action=", val, "; "}, "")
	}

	return ""
}

func withErrorTag(ctx context.Context) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("error", true)
	}
}
