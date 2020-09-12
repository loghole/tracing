package tracehttp

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type key string

const (
	NextRoundTrip key = "next"
)

type Transport struct {
	tracer opentracing.Tracer
	base   http.RoundTripper
}

func NewTransport(tracer opentracing.Tracer, roundTripper http.RoundTripper) *Transport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}

	return &Transport{
		tracer: tracer,
		base:   roundTripper,
	}
}

// RoundTrip implements the RoundTripper interface.
func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var parentCtx opentracing.SpanContext

	if parent := opentracing.SpanFromContext(req.Context()); parent != nil {
		parentCtx = parent.Context()
	}

	span := t.tracer.StartSpan(buildSpanName(req), opentracing.ChildOf(parentCtx))
	defer span.Finish()

	ext.Component.Set(span, "http")
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())

	err = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		log.Printf("[error] inject headers failed: %v", err)
	}

	resp, err = t.base.RoundTrip(req.WithContext(opentracing.ContextWithSpan(req.Context(), span)))
	if err != nil {
		ext.Error.Set(span, true)

		return resp, err
	}

	ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))

	return resp, nil
}

func (t *Transport) checkNextRoundTrip(req *http.Request) (*http.Request, bool) {
	if val := req.Context().Value(NextRoundTrip); val != nil {
		return req, true
	}

	return req.WithContext(context.WithValue(req.Context(), NextRoundTrip, true)), false
}

func buildSpanName(r *http.Request) string {
	return strings.Join([]string{"HTTP", r.Method, r.URL.String()}, " ")
}
