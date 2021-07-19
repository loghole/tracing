package logtracer

import (
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

type SpanContext interface {
	opentracing.SpanContext
	TraceID() jaeger.TraceID
	SpanID() jaeger.SpanID
}
