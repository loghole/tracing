package mocks

import (
	"context"
	"encoding/binary"

	"go.opentelemetry.io/otel/trace"
)

// NewContextWithMockSpan is test helper.
func NewContextWithMockSpan(ctx context.Context, traceID, spanID uint64) context.Context {
	var (
		spanContext          trace.SpanContext
		opentelemetryTraceID trace.TraceID
		opentelemetrySpanID  trace.SpanID
	)

	binary.LittleEndian.PutUint64(opentelemetryTraceID[:8], traceID)
	binary.LittleEndian.PutUint64(opentelemetrySpanID[:8], spanID)

	spanContext = spanContext.WithTraceID(opentelemetryTraceID)
	spanContext = spanContext.WithSpanID(opentelemetrySpanID)

	return trace.ContextWithSpanContext(ctx, spanContext)
}
