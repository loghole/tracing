package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// SpanFromContext returns the current Span from ctx.
// If no Span is currently set in ctx an implementation of a Span that performs no operations is returned.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// InjectMap set tracecontext from the Context into the map[string]string carrier.
func InjectMap(ctx context.Context, carrier map[string]string) {
	new(propagation.TraceContext).Inject(ctx, propagation.MapCarrier(carrier))
}

// InjectHeaders set tracecontext from the Context into the http.Header carrier.
func InjectHeaders(ctx context.Context, carrier http.Header) {
	new(propagation.TraceContext).Inject(ctx, propagation.HeaderCarrier(carrier))
}
