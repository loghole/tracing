package logtracer

import (
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
)

type LogTracer struct {
	rnd *rand.Rand
}

func NewLogTracer() *LogTracer {
	return &LogTracer{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))} // nolint:gosec // need math rnd
}

// StartSpan belongs to the Tracer interface.
func (t *LogTracer) StartSpan(_ string, opts ...opentracing.StartSpanOption) opentracing.Span {
	s := &opentracing.StartSpanOptions{}

	for _, opt := range opts {
		opt.Apply(s)

		if len(s.References) > 0 {
			break
		}
	}

	for _, ref := range s.References {
		if val, ok := ref.ReferencedContext.(*LogSpanContext); ok {
			return &LogSpan{traceID: val.span.traceID, spanID: t.generateSpanID(), tracer: t}
		}
	}

	return &LogSpan{traceID: t.generateTraceID(), spanID: t.generateSpanID(), tracer: t}
}

// Inject belongs to the Tracer interface.
func (t *LogTracer) Inject(sp opentracing.SpanContext, _, carrier interface{}) error {
	m, ok := carrier.(opentracing.TextMapCarrier)
	if !ok {
		return nil
	}

	if val, ok := sp.(*LogSpanContext); ok {
		m[jaeger.TraceContextHeaderName] = val.span.spanID.String()
	}

	return nil
}

// Extract belongs to the Tracer interface.
func (t LogTracer) Extract(_, _ interface{}) (opentracing.SpanContext, error) {
	return nil, opentracing.ErrSpanContextNotFound
}

func (t *LogTracer) generateTraceID() jaeger.TraceID {
	return jaeger.TraceID{Low: t.rnd.Uint64()}
}

func (t *LogTracer) generateSpanID() jaeger.SpanID {
	return jaeger.SpanID(t.rnd.Uint64())
}

type LogSpanContext struct {
	span LogSpan
}

func (n *LogSpanContext) ForeachBaggageItem(_ func(k, v string) bool) {}

func (n *LogSpanContext) TraceID() jaeger.TraceID {
	return n.span.traceID
}

func (n *LogSpanContext) SpanID() jaeger.SpanID {
	return n.span.spanID
}

type LogSpan struct {
	tracer  *LogTracer
	traceID jaeger.TraceID
	spanID  jaeger.SpanID
}

func (s *LogSpan) Context() opentracing.SpanContext                { return &LogSpanContext{span: *s} }
func (s *LogSpan) SetBaggageItem(_, _ string) opentracing.Span     { return s }
func (s *LogSpan) BaggageItem(_ string) string                     { return "" }
func (s *LogSpan) SetTag(_ string, _ interface{}) opentracing.Span { return s }
func (s *LogSpan) LogFields(_ ...log.Field)                        {}
func (s *LogSpan) LogKV(_ ...interface{})                          {}
func (s *LogSpan) Finish()                                         {}
func (s *LogSpan) FinishWithOptions(_ opentracing.FinishOptions)   {}
func (s *LogSpan) SetOperationName(_ string) opentracing.Span      { return s }
func (s *LogSpan) Tracer() opentracing.Tracer                      { return s.tracer }
func (s *LogSpan) LogEvent(_ string)                               {}
func (s *LogSpan) LogEventWithPayload(_ string, _ interface{})     {}
func (s *LogSpan) Log(_ opentracing.LogData)                       {}
