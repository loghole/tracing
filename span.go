package tracing

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	baseSkipCallers = 3
)

type Span struct {
	tracer opentracing.Tracer
	span   opentracing.Span

	finished uint32
}

func ChildSpan(ctx *context.Context) (s *Span) { // nolint:gocritic
	s = &Span{}

	if span := opentracing.SpanFromContext(*ctx); span != nil {
		s.span = span.Tracer().StartSpan(
			callerName(),
			opentracing.ChildOf(span.Context()),
			opentracing.StartTime(time.Now()),
		)

		// Overriding the original context
		*ctx = opentracing.ContextWithSpan(*ctx, s.span)
	}

	return s
}

func FollowsSpan(ctx *context.Context) (s *Span) { // nolint:gocritic
	s = &Span{}

	if span := opentracing.SpanFromContext(*ctx); span != nil {
		s.span = span.Tracer().StartSpan(
			callerName(),
			opentracing.FollowsFrom(span.Context()),
			opentracing.StartTime(time.Now()),
		)

		// Overriding the original context
		*ctx = opentracing.ContextWithSpan(context.Background(), s.span)
	}

	return s
}

func (s *Span) Finish() {
	if !atomic.CompareAndSwapUint32(&s.finished, 0, 1) {
		warnf("%s finish finished span", callerLine())
	}

	if s.span != nil {
		s.span.Finish()
	}
}

func (s *Span) FinishWithOptions(opts opentracing.FinishOptions) {
	if !atomic.CompareAndSwapUint32(&s.finished, 0, 1) {
		warnf("%s finish finished span", callerLine())
	}

	if s.span != nil {
		s.span.FinishWithOptions(opts)
	}
}

func (s *Span) Context() opentracing.SpanContext {
	if s.span != nil {
		return s.span.Context()
	}

	return nil
}

func (s *Span) SetOperationName(operationName string) opentracing.Span {
	if s.span != nil {
		return s.span.SetOperationName(operationName)
	}

	return s
}

func (s *Span) SetTag(key string, value interface{}) opentracing.Span {
	if s.span != nil {
		return s.span.SetTag(key, value)
	}

	return s
}

func (s *Span) LogFields(fields ...log.Field) {
	if s.span != nil {
		s.span.LogFields(fields...)
	}
}

func (s *Span) LogKV(alternatingKeyValues ...interface{}) {
	if s.span != nil {
		s.span.LogKV(alternatingKeyValues...)
	}
}

func (s *Span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	if s.span != nil {
		return s.span.SetBaggageItem(restrictedKey, value)
	}

	return s
}

func (s *Span) BaggageItem(restrictedKey string) string {
	if s.span != nil {
		return s.span.BaggageItem(restrictedKey)
	}

	return ""
}

func (s *Span) Tracer() opentracing.Tracer {
	if s.tracer != nil {
		return s.tracer
	}

	return nil
}

// Deprecated: use LogFields or LogKV
func (s *Span) LogEvent(event string) {
	if s.span != nil {
		s.span.LogFields(log.String(event, ""))
	}
}

// Deprecated: use LogFields or LogKV
func (s *Span) LogEventWithPayload(event string, payload interface{}) {
	if s.span != nil {
		s.span.LogKV(event, payload)
	}
}

// Deprecated: use LogFields or LogKV
func (s *Span) Log(data opentracing.LogData) {
	if s.span != nil {
		s.span.LogKV(data.Event, data.Payload)
	}
}

func callerName() string {
	var pc [1]uintptr

	runtime.Callers(baseSkipCallers, pc[:])

	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return "unknown"
	}

	list := strings.Split(f.Name(), "/")

	return list[len(list)-1]
}

func callerLine() string {
	var pc [1]uintptr

	runtime.Callers(baseSkipCallers, pc[:])

	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return "unknown"
	}

	file, line := f.FileLine(pc[0])

	return fmt.Sprintf("%s:%d", file, line)
}
