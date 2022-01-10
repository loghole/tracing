package tracing

import (
	"context"
	"encoding/json"
	"io"
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

func InjectBinary(ctx context.Context, carrier io.Writer) {
	if span := SpanFromContext(ctx); span != nil {
		data := make(map[string]string)

		new(propagation.TraceContext).Inject(ctx, propagation.MapCarrier(data))

		_ = json.NewEncoder(carrier).Encode(data)
	}
}

func InjectHeaders(ctx context.Context, carrier http.Header) {
	new(propagation.TraceContext).Inject(ctx, propagation.HeaderCarrier(carrier))
}
