package tracing

import (
	"bytes"
	"context"
	"net/http"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
)

const (
	_skipCallers = 3
)

type Span struct {
	span opentracing.Span
}

func ChildSpan(ctx *context.Context) (tracer Span) { // nolint:gocritic
	if span := opentracing.SpanFromContext(*ctx); span != nil {
		tracer.span = span.Tracer().StartSpan(callerName(), opentracing.ChildOf(span.Context()))

		// Переопределяем исходный контекст
		*ctx = opentracing.ContextWithSpan(*ctx, tracer.span)
	}

	return tracer
}

func FollowsSpan(ctx *context.Context) (tracer Span) { // nolint:gocritic
	if span := opentracing.SpanFromContext(*ctx); span != nil {
		tracer.span = span.Tracer().StartSpan(callerName(), opentracing.FollowsFrom(span.Context()))

		// Переопределяем исходный контекст
		*ctx = opentracing.ContextWithSpan(*ctx, tracer.span)
	}

	return tracer
}

func (s Span) WithTag(key string, val interface{}) Span {
	if s.span != nil {
		s.span.SetTag(key, val)
	}

	return s
}

func (s Span) Finish() {
	if s.span != nil {
		s.span.Finish()
	}
}

func (s Span) Context(ctx context.Context) context.Context {
	return opentracing.ContextWithSpan(ctx, s.span)
}

func (s Span) GetSpanContext() opentracing.SpanContext {
	if s.span != nil {
		return s.span.Context()
	}

	return nil
}

func InjectMap(ctx context.Context) map[string]string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		carrier := opentracing.TextMapCarrier{}

		err := span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier)
		if err == nil {
			return carrier
		}
	}

	return nil
}

func InjectBinary(ctx context.Context) []byte {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		carrier := bytes.NewBuffer([]byte{})

		err := span.Tracer().Inject(span.Context(), opentracing.Binary, carrier)
		if err == nil {
			return carrier.Bytes()
		}
	}

	return nil
}

func InjectHeaders(ctx context.Context) http.Header {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		carrier := opentracing.HTTPHeadersCarrier{}

		err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err == nil {
			return http.Header(carrier)
		}
	}

	return nil
}

func callerName() string {
	var pc [1]uintptr

	runtime.Callers(_skipCallers, pc[:])
	f := runtime.FuncForPC(pc[0])

	list := strings.Split(f.Name(), "/")

	return list[len(list)-1]
}
