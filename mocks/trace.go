package mocks

import (
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func NewTracerWithRecorder() (trace.Tracer, *tracetest.SpanRecorder) {
	var (
		recorder = tracetest.NewSpanRecorder()
		tracer   = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(recorder)).Tracer("")
	)

	return tracer, recorder
}
