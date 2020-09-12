package tracehttp

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func NewClient(tracer opentracing.Tracer, client *http.Client) *http.Client {
	client.Transport = NewTransport(tracer, client.Transport)

	return client
}
