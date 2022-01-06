package main

import (
	"context"
	"log"
	"time"

	"github.com/loghole/tracing/internal/spanprocessor"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service     = "trace-demo"
	environment = "dev"
	id          = 1
)

func main() {
	const uri = "https://tracing.3l8.ru/api/input"

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(uri)))
	if err != nil {
		panic(err)
	}

	sampledProcessor := spanprocessor.NewSampled(
		tracesdk.NewBatchSpanProcessor(exp, tracesdk.WithMaxExportBatchSize(1)),
		tracesdk.AlwaysSample(),
	)

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(sampledProcessor),

		// Always be sure to batch in production.
		// tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		time.Sleep(time.Second * 10)

		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-main")

	log.Println("start foo trace")
	ctx, span := tr.Start(ctx, "foo2")
	defer span.End()

	time.Sleep(time.Millisecond * 250)

	span.AddEvent("test1", trace.WithTimestamp(time.Now()), trace.WithAttributes(attribute.KeyValue{"atr1", attribute.StringValue("val1")}))

	time.Sleep(time.Millisecond * 250)

	span.AddEvent("test2", trace.WithTimestamp(time.Now()))

	bar(ctx, tp, "bar1")

	bar(ctx, tp, "bar2")
	bar2(ctx, tp, "bar3")

	log.Println("end")
}

/*
	Указать время спана

	ctx, span := tr.Start(ctx, "foo", trace.WithTimestamp(time.Now().Add(-5*time.Minute)))
	defer span.End()
*/
func bar(ctx context.Context, tp *tracesdk.TracerProvider, name string) {
	time.Sleep(time.Second)
	// Use the global TracerProvider.
	//tr := otel.Tracer("component-bar")
	ctx, span := tp.Tracer("component-bar").Start(ctx, name)
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	time.Sleep(time.Millisecond * 250)

	//setErrorTag(ctx)

	span.AddEvent("test3", trace.WithTimestamp(time.Now()))

	// Do bar...
}

func bar2(ctx context.Context, tp *tracesdk.TracerProvider, name string) {
	time.Sleep(time.Second)
	// Use the global TracerProvider.
	//tr := otel.Tracer("component-bar")
	ctx, span := tp.Tracer("component-bar").Start(ctx, name)
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer func() {
		go func() {
			time.Sleep(time.Second * 5)
			span.End()
		}()
	}()

	time.Sleep(time.Millisecond * 250)

	//setErrorTag(ctx)

	span.AddEvent("test3", trace.WithTimestamp(time.Now()))

	// Do bar...
}

func setErrorTag(ctx context.Context) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.KeyValue{"error", attribute.BoolValue(true)})
}
