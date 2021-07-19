package test

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
)

type MockSpan struct {
	traceID jaeger.TraceID
	spanID  jaeger.SpanID
}

func NewMockSpan(traceID, spanID uint64) *MockSpan {
	return &MockSpan{
		traceID: jaeger.TraceID{Low: traceID},
		spanID:  jaeger.SpanID(spanID),
	}
}

func NewContextWithMockSpan(ctx context.Context, traceID, spanID uint64) context.Context {
	return opentracing.ContextWithSpan(ctx, NewMockSpan(traceID, spanID))
}

func (s *MockSpan) Context() opentracing.SpanContext                { return &MockSpanContext{span: *s} }
func (s *MockSpan) SetBaggageItem(_, _ string) opentracing.Span     { return s }
func (s *MockSpan) BaggageItem(_ string) string                     { return "" }
func (s *MockSpan) SetTag(_ string, _ interface{}) opentracing.Span { return s }
func (s *MockSpan) LogFields(_ ...log.Field)                        {}
func (s *MockSpan) LogKV(_ ...interface{})                          {}
func (s *MockSpan) Finish()                                         {}
func (s *MockSpan) FinishWithOptions(_ opentracing.FinishOptions)   {}
func (s *MockSpan) SetOperationName(_ string) opentracing.Span      { return s }
func (s *MockSpan) Tracer() opentracing.Tracer                      { return nil }
func (s *MockSpan) LogEvent(_ string)                               {}
func (s *MockSpan) LogEventWithPayload(_ string, _ interface{})     {}
func (s *MockSpan) Log(_ opentracing.LogData)                       {}

type MockSpanContext struct {
	span MockSpan
}

func (n *MockSpanContext) ForeachBaggageItem(_ func(k, v string) bool) {}

func (n *MockSpanContext) TraceID() jaeger.TraceID {
	return n.span.traceID
}

func (n *MockSpanContext) SpanID() jaeger.SpanID {
	return n.span.spanID
}
