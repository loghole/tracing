package tracehttp

import (
	"log"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
	var span opentracing.Span

	if parent := opentracing.SpanFromContext(req.Context()); parent != nil {
		span = t.tracer.StartSpan(buildSpanName(req), opentracing.ChildOf(parent.Context()))
		defer span.Finish()

		ext.SpanKindRPCClient.Set(span)
		ext.HTTPMethod.Set(span, req.Method)
		ext.HTTPUrl.Set(span, req.URL.String())

		err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			log.Printf("[error] inject headers failed: %v", err)
		}
	}

	resp, err = t.base.RoundTrip(req)
	if err != nil {
		if span != nil {
			ext.Error.Set(span, true)
		}

		return resp, err
	}

	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	}

	return resp, nil
}

func buildSpanName(r *http.Request) string {
	return strings.Join([]string{"HTTP", r.Method, r.URL.String()}, " ")
}
