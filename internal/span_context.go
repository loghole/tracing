package internal

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

type SpanContext interface {
	opentracing.SpanContext
	TraceID() jaeger.TraceID
	SpanID() jaeger.SpanID
}
