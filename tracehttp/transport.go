package tracehttp

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/loghole/tracing"
)

type Transport struct {
	tracer   *tracing.Tracer
	base     http.RoundTripper
	extended http.RoundTripper
}

func NewTransport(tracer *tracing.Tracer, roundTripper http.RoundTripper) *Transport {
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
	if t.extended != nil {
		return t.extended.RoundTrip(req)
	}

	ctx, span := t.tracer.NewSpan().WithName(t.defaultNameFunc(req)).StartWithContext(req.Context())
	defer span.End()

	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
		attribute.String("component", ComponentName),
	)

	tracing.InjectHeaders(ctx, req.Header)

	resp, err = t.base.RoundTrip(req.WithContext(ctx))
	if err != nil {
		span.SetAttributes(attribute.Bool("error", true))

		return resp, err
	}

	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(resp.StatusCode))

	return resp, nil
}

func (t *Transport) defaultNameFunc(req *http.Request) string {
	return strings.Join([]string{"HTTP", req.Method, req.RequestURI}, " ")
}
