package tracehttp

import (
	"net/http"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/loghole/tracing"
)

type Transport struct {
	tracer trace.Tracer
	base   http.RoundTripper
}

func NewTransport(tracer trace.Tracer, roundTripper http.RoundTripper) *Transport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}

	transport := &Transport{
		tracer: tracer,
		base:   roundTripper,
	}

	return transport
}

// RoundTrip implements the RoundTripper interface.
func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx, span := t.tracer.Start(req.Context(), defaultNameFunc(req), trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
		semconv.HTTPSchemeKey.String(req.URL.Scheme),
		semconv.HTTPRequestContentLengthKey.Int64(req.ContentLength),
	)

	tracing.InjectHeaders(ctx, req.Header)

	resp, err = t.base.RoundTrip(req.WithContext(ctx))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error")

		return resp, err
	}

	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(resp.StatusCode))

	return resp, nil
}

func defaultNameFunc(req *http.Request) string {
	return "HTTP " + req.Method + " " + req.RequestURI
}
