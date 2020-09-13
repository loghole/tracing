package tracehttp

import (
	"log"
	"net/http"
	"strings"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Transport struct {
	tracer   opentracing.Tracer
	base     http.RoundTripper
	extended http.RoundTripper
}

func NewTransport(tracer opentracing.Tracer, roundTripper http.RoundTripper, extended bool) *Transport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}

	transport := &Transport{
		tracer: tracer,
		base:   roundTripper,
	}

	if extended {
		transport.extended = &nethttp.Transport{RoundTripper: roundTripper}
	}

	return transport
}

// RoundTrip implements the RoundTripper interface.
func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if t.extended != nil {
		return t.extended.RoundTrip(req)
	}

	var parentCtx opentracing.SpanContext

	if parent := opentracing.SpanFromContext(req.Context()); parent != nil {
		parentCtx = parent.Context()
	}

	span := t.tracer.StartSpan(t.defaultNameFunc(req), opentracing.ChildOf(parentCtx))
	defer span.Finish()

	ext.Component.Set(span, ComponentName)
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

func (t *Transport) defaultNameFunc(req *http.Request) string {
	return strings.Join([]string{"HTTP", req.Method, req.RequestURI}, " ")
}
