package tracing

import (
	"context"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// LogKV is a concise, readable way to record key:value logging data about
// a Span, though unfortunately this also makes it less efficient and less
// type-safe than LogFields(). Here's an example:
//
//    span.LogKV(
//        "event", "soft error",
//        "type", "cache timeout",
//        "waited.millis", 1500)
//
// For LogKV (as opposed to LogFields()), the parameters must appear as
// key-value pairs, like
//
//    span.LogKV(key1, val1, key2, val2, key3, val3, ...)
//
// The keys must all be strings. The values may be strings, numeric types,
// bools, Go error instances, or arbitrary structs.
//
// (Note to implementors: consider the log.InterleavedKVToFields() helper).
func LogKV(ctx context.Context, key string, value interface{}) {
	if span := SpanFromContext(ctx); span != nil {
		span.LogKV(key, value)
	}
}

func InjectMap(ctx context.Context, carrier map[string]string) error {
	if span := SpanFromContext(ctx); span != nil {
		err := span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(carrier))
		if err != nil {
			return err
		}
	}

	return nil
}

func InjectBinary(ctx context.Context, carrier io.Writer) error {
	if span := SpanFromContext(ctx); span != nil {
		err := span.Tracer().Inject(span.Context(), opentracing.Binary, carrier)
		if err != nil {
			return err
		}
	}

	return nil
}

func InjectHeaders(ctx context.Context, carrier http.Header) error {
	if span := SpanFromContext(ctx); span != nil {
		err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(carrier))
		if err != nil {
			return err
		}
	}

	return nil
}
