package spanprocessor

import (
	"context"
	"strings"
	"sync"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type spanWrapper struct {
	span    tracesdk.ReadWriteSpan
	context context.Context
}

type traceWrapper struct {
	mu    sync.Mutex
	spans map[trace.SpanID]spanWrapper

	parentSpanID trace.SpanID
	isFinished   bool
	_hasError    bool
}

func (t *traceWrapper) storeSpan(ctx context.Context, span tracesdk.ReadWriteSpan) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.spans[span.SpanContext().SpanID()] = spanWrapper{
		context: ctx,
		span:    span,
	}
}

func (t *traceWrapper) extractFinishedSpans() []spanWrapper {
	t.mu.Lock()
	defer t.mu.Unlock()

	result := make([]spanWrapper, 0, len(t.spans))

	for key, span := range t.spans {
		if span.span.EndTime().IsZero() {
			continue
		}

		result = append(result, span)

		delete(t.spans, key)
	}

	return result
}

func (t *traceWrapper) isParent(spanID trace.SpanID) bool {
	return t.parentSpanID == spanID
}

func (t *traceWrapper) hasError() bool {
	if t._hasError {
		return true
	}

	for _, span := range t.spans {
		for _, attr := range span.span.Attributes() {
			if strings.EqualFold(string(attr.Key), "error") {
				t._hasError = true

				return true
			}
		}
	}

	return false
}
