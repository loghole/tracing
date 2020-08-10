package tracing

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/opentracing/opentracing-go"
)

const (
	_skipCallers = 3
)

type Span struct {
	span opentracing.Span
	once sync.Once
	done uint32
}

func ChildSpan(ctx *context.Context) (s *Span) { // nolint:gocritic
	s = &Span{}

	if span := opentracing.SpanFromContext(*ctx); span != nil {
		s.span = span.Tracer().StartSpan(callerName(), opentracing.ChildOf(span.Context()))

		// Overriding the original context
		*ctx = opentracing.ContextWithSpan(*ctx, s.span)
	}

	return s
}

func FollowsSpan(ctx *context.Context) (s *Span) { // nolint:gocritic
	s = &Span{}

	if span := opentracing.SpanFromContext(*ctx); span != nil {
		s.span = span.Tracer().StartSpan(callerName(), opentracing.FollowsFrom(span.Context()))

		// Overriding the original context
		*ctx = opentracing.ContextWithSpan(*ctx, s.span)
	}

	return s
}

func (s *Span) WithTag(key string, val interface{}) *Span {
	if s.span != nil {
		s.span.SetTag(key, val)
	}

	return s
}

func (s *Span) Finish() {
	if atomic.LoadUint32(&s.done) > 0 {
		log.New(os.Stdout, "tracing: ", log.Ldate|log.Llongfile).Println("[warn] finish finished span")
	}

	if s.span != nil {
		s.once.Do(s.span.Finish)
	}

	atomic.AddUint32(&s.done, 1)
}

func (s *Span) Context(ctx context.Context) context.Context {
	return opentracing.ContextWithSpan(ctx, s.span)
}

func (s *Span) GetSpanContext() opentracing.SpanContext {
	if s.span != nil {
		return s.span.Context()
	}

	return nil
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

	runtime.Callers(_skipCallers, pc[:])

	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return "unknown"
	}

	list := strings.Split(f.Name(), "/")

	return list[len(list)-1]
}
