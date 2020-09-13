package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gadavy/tracing"
	"github.com/gadavy/tracing/tracelog"
)

const (
	jaegerURL  = "127.0.0.1:6831"
)

func main() {
	dev, err := zap.NewProduction(zap.AddStacktrace(zap.FatalLevel))
	if err != nil {
		panic(err)
	}

	logger := tracelog.NewTraceLogger(dev.Sugar())

	fmt.Println("============== Base example =============")
	exampleBase := NewBaseExample(logger)

	fmt.Println("\n============== #1 ======================")
	exampleBase.CreateSpanBase()

	fmt.Println("\n============== #2 ======================")
	exampleBase.CreateSpanWithHTTP(nil, &http.Request{})
}

type BaseExample struct {
	tracer *tracing.Tracer
	logger tracelog.Logger
}

func NewBaseExample(logger tracelog.Logger) *BaseExample {
	tracer, err := tracing.NewTracer(tracing.DefaultConfiguration("base_example", jaegerURL))
	if err != nil {
		panic(err)
	}

	return &BaseExample{
		tracer: tracer,
		logger: logger,
	}
}

func (e *BaseExample) CreateSpanBase() {
	span, ctx := e.tracer.NewSpan().
		WithName("ExampleCreateSpanBase").
		BuildWithContext(context.Background())
	defer span.Finish()

	e.logger.Info(ctx, "ExampleCreateSpanBase info message")

	e.ChildSpanFromContext(ctx)
}

func (e *BaseExample) CreateSpanWithHTTP(w http.ResponseWriter, r *http.Request) {
	span, ctx := e.tracer.NewSpan().
		WithName("ExampleCreateSpanWithHTTP").
		ExtractHeaders(r.Header).
		BuildWithContext(context.Background())
	defer span.Finish()

	e.logger.Info(ctx, "ExampleCreateSpanWithHTTP info message")

	e.ChildSpanFromContext(ctx)
}

func (e *BaseExample) ChildSpanFromContext(ctx context.Context) {
	defer tracing.ChildSpan(&ctx).SetTag("key", "val").Finish()

	time.Sleep(time.Second)

	// Trace logger example.
	e.logger.Debug(ctx, "ExampleChildSpanFromContext debug message")
	e.logger.Info(ctx, "ExampleChildSpanFromContext info message")
	e.logger.Warn(ctx, "ExampleChildSpanFromContext warn message")

	// Trace logger error set error tag to span.
	e.logger.Error(ctx, "ExampleChildSpanFromContext error message")
}

func (e *BaseExample) Close() error {
	return e.tracer.Close()
}

