package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/gadavy/tracing"
)

var (
	service   = ""
	jaegerURL = ""
)

func main() {
	dev, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger := tracing.DefaultTraceLogger(dev.Sugar())

	tracer, err := tracing.NewTracer(tracing.DefaultConfiguration(service, jaegerURL))
	if err != nil {
		panic(err)
	}

	example := NewExample(tracer, logger)

	example.ExampleCreateSpanBase()
}

type Example struct {
	tracer *tracing.Tracer
	logger *tracing.TraceLogger
}

func NewExample(tracer *tracing.Tracer, logger *tracing.TraceLogger) *Example {
	return &Example{
		tracer: tracer,
		logger: logger,
	}
}

func (e *Example) ExampleCreateSpanBase() {
	span, ctx := e.tracer.NewSpan().
		WithName("ExampleCreateSpanBase").
		BuildWithContext(context.Background())
	defer span.Finish()

	e.logger.Info(ctx, "ExampleCreateSpanBase message 1")

	e.ExampleChildSpanFromContext(ctx)
}

func (e *Example) ExampleChildSpanFromContext(ctx context.Context) {
	defer tracing.ChildSpan(&ctx).WithTag("key", "val").Finish()

	e.logger.Info(ctx, "ExampleChildSpanFromContext message 1")
}
