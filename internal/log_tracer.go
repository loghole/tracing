package internal

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
)

type LogTracer struct {
	rnd *rand.Rand
}

func NewLogTracer() *LogTracer {
	return &LogTracer{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
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
		if val, ok := ref.ReferencedContext.(*logSpanContext); ok {
			return &logSpan{traceID: val.traceID, tracer: t}
		}
	}

	return &logSpan{traceID: t.generateTraceID(), tracer: t}
}

// Inject belongs to the Tracer interface.
func (t *LogTracer) Inject(sp opentracing.SpanContext, _, carrier interface{}) error {
	m, ok := carrier.(opentracing.TextMapCarrier)
	if !ok {
		return nil
	}

	if val, ok := sp.(*logSpanContext); ok {
		m[jaeger.TraceContextHeaderName] = val.traceID
	}

	return nil
}

func (t *LogTracer) generateTraceID() string {
	return strconv.FormatInt(t.rnd.Int63(), 16)
}

// Extract belongs to the Tracer interface.
func (t LogTracer) Extract(_, _ interface{}) (opentracing.SpanContext, error) {
	return nil, opentracing.ErrSpanContextNotFound
}

type logSpanContext struct {
	traceID string
}

// logSpanContext:
func (n *logSpanContext) ForeachBaggageItem(_ func(k, v string) bool) {}

type logSpan struct {
	tracer  *LogTracer
	traceID string
}

func (n *logSpan) Context() opentracing.SpanContext                { return &logSpanContext{traceID: n.traceID} }
func (n *logSpan) SetBaggageItem(_, _ string) opentracing.Span     { return n }
func (n *logSpan) BaggageItem(_ string) string                     { return "" }
func (n *logSpan) SetTag(_ string, _ interface{}) opentracing.Span { return n }
func (n *logSpan) LogFields(_ ...log.Field)                        {}
func (n *logSpan) LogKV(_ ...interface{})                          {}
func (n *logSpan) Finish()                                         {}
func (n *logSpan) FinishWithOptions(_ opentracing.FinishOptions)   {}
func (n *logSpan) SetOperationName(_ string) opentracing.Span      { return n }
func (n *logSpan) Tracer() opentracing.Tracer                      { return n.tracer }
func (n *logSpan) LogEvent(_ string)                               {}
func (n *logSpan) LogEventWithPayload(_ string, _ interface{})     {}
func (n *logSpan) Log(_ opentracing.LogData)                       {}
