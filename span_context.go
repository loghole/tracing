package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func InjectMap(ctx context.Context, carrier map[string]string) {
	new(propagation.TraceContext).Inject(ctx, propagation.MapCarrier(carrier))
}

func InjectHeaders(ctx context.Context, carrier http.Header) {
	new(propagation.TraceContext).Inject(ctx, propagation.HeaderCarrier(carrier))
}
