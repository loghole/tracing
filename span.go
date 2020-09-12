package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

	s.span.Finish()
}

func (s *Span) FinishWithOptions(opts opentracing.FinishOptions) {
	if !atomic.CompareAndSwapUint32(&s.finished, 0, 1) {
		warnf("%s finish finished span", callerLine())
	}

	s.span.FinishWithOptions(opts)
}

func (s *Span) Context() opentracing.SpanContext {
	return s.span.Context()
}

func (s *Span) SetOperationName(operationName string) opentracing.Span {
	return s.span.SetOperationName(operationName)
}

func (s *Span) SetTag(key string, value interface{}) opentracing.Span {
	return s.span.SetTag(key, value)
}

func (s *Span) LogFields(fields ...log.Field) {
	s.span.LogFields(fields...)
}

func (s *Span) LogKV(alternatingKeyValues ...interface{}) {
	s.span.LogKV(alternatingKeyValues...)
}

func (s *Span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	return s.span.SetBaggageItem(restrictedKey, value)
}

func (s *Span) BaggageItem(restrictedKey string) string {
	return s.span.BaggageItem(restrictedKey)
}

func (s *Span) Tracer() opentracing.Tracer {
	return s.tracer
}

// Deprecated: use LogFields or LogKV
func (s *Span) LogEvent(event string) {
	s.span.LogFields(log.String(event, ""))
}

// Deprecated: use LogFields or LogKV
func (s *Span) LogEventWithPayload(event string, payload interface{}) {
	s.span.LogKV(event, payload)
}

// Deprecated: use LogFields or LogKV
func (s *Span) Log(data opentracing.LogData) {
	s.span.LogKV(data.Event, data.Payload)
}

func InjectMap(ctx context.Context, carrier map[string]string) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		err := span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(carrier))
		if err != nil {
			return err
		}
	}

	return nil
}

func InjectBinary(ctx context.Context, carrier io.Writer) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		err := span.Tracer().Inject(span.Context(), opentracing.Binary, carrier)
		if err != nil {
			return err
		}
	}

	return nil
}

func InjectHeaders(ctx context.Context, carrier http.Header) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(carrier))
		if err != nil {
			return err
		}
	}

	return nil
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
