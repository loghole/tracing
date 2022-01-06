package main

import (
	"context"
	"log"
	"time"

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

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(&OurSpanProcessor{
			root: tracesdk.NewBatchSpanProcessor(exp),
		}),

		// Always be sure to batch in production.
		// tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
		tracesdk.WithSampler(&OurSampler{}),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-main")

	log.Println("start foo trace")
	ctx, span := tr.Start(ctx, "foo")
	defer span.End()

	time.Sleep(time.Millisecond * 250)

	span.AddEvent("test1", trace.WithTimestamp(time.Now()), trace.WithAttributes(attribute.KeyValue{"atr1", attribute.StringValue("val1")}))

	time.Sleep(time.Millisecond * 250)

	span.AddEvent("test2", trace.WithTimestamp(time.Now()))
	/*
		ctx = trace.ContextWithSpanContext(
			ctx,
			span.SpanContext().WithTraceFlags(span.SpanContext().TraceFlags().WithSampled(false)),
		)*/

	bar(ctx, tp)

	log.Println("end")
}

/*
	Указать время спана

	ctx, span := tr.Start(ctx, "foo", trace.WithTimestamp(time.Now().Add(-5*time.Minute)))
	defer span.End()
*/

func bar(ctx context.Context, tp *tracesdk.TracerProvider) {
	time.Sleep(time.Second)
	// Use the global TracerProvider.
	//tr := otel.Tracer("component-bar")
	ctx, span := tp.Tracer("component-bar").Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	time.Sleep(time.Millisecond * 250)

	setErrorTag(ctx)

	span.AddEvent("test3", trace.WithTimestamp(time.Now()))

	// Do bar...
}

func setErrorTag(ctx context.Context) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.KeyValue{"error", attribute.BoolValue(true)})

}

type Sampler interface {
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// ShouldSample returns a SamplingResult based on a decision made from the
	// passed parameters.
	ShouldSample(parameters tracesdk.SamplingParameters) tracesdk.SamplingResult
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// Description returns information describing the Sampler.
	Description() string
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.
}

type OurSampler struct {
}

func (s *OurSampler) ShouldSample(parameters tracesdk.SamplingParameters) tracesdk.SamplingResult {
	log.Println("ShouldSample")

	return tracesdk.SamplingResult{
		Decision: tracesdk.RecordAndSample,
		Attributes: []attribute.KeyValue{
			attribute.KeyValue{"IsRecording", attribute.BoolValue(true)},
			attribute.KeyValue{"Sampled", attribute.BoolValue(true)},
		},
	}
}

func (s *OurSampler) Description() string {
	log.Println("Description")

	return ""
}

/*

// ShouldSample implements the Sampler's ShouldSample().
//nolint:gocritic
func (rSampler *RemotelyControlledSampler) ShouldSample(param otelTrace.SamplingParameters) otelTrace.SamplingResult {
	sampler := rSampler.getSamplerForOperation(param.Name)
	isSampled, tags := sampler.IsSampled(param.TraceID)
	tags = append(tags, label.Bool("IsRecording", true), label.Bool("Sampled", isSampled))

	return otelTrace.SamplingResult{
		Decision:   otelTrace.RecordAndSample,
		Attributes: tags,
	}
}
*/

type SpanProcessor interface {

	// OnStart is called when a span is started. It is called synchronously
	// and should not block.
	OnStart(parent context.Context, s tracesdk.ReadWriteSpan)

	// OnEnd is called when span is finished. It is called synchronously and
	// hence not block.
	OnEnd(s tracesdk.ReadOnlySpan)

	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the processor should be done in this call.
	//
	// Calls to OnStart, OnEnd, or ForceFlush after this has been called
	// should be ignored.
	//
	// All timeouts and cancellations contained in ctx must be honored, this
	// should not block indefinitely.
	Shutdown(ctx context.Context) error

	// ForceFlush exports all ended spans to the configured Exporter that have not yet
	// been exported.  It should only be called when absolutely necessary, such as when
	// using a FaaS provider that may suspend the process after an invocation, but before
	// the Processor can export the completed spans.
	ForceFlush(ctx context.Context) error
}

type OurSpanProcessor struct {
	root tracesdk.SpanProcessor
}

func (p *OurSpanProcessor) OnStart(parent context.Context, s tracesdk.ReadWriteSpan) {
	log.Printf("on start: %+v", s)

	p.root.OnStart(parent, s)
}

func (p *OurSpanProcessor) OnEnd(s tracesdk.ReadOnlySpan) {
	log.Printf("on end: %+v", s)

	p.root.OnEnd(s)
}

func (p *OurSpanProcessor) Shutdown(ctx context.Context) error {
	log.Println("on shutdown")

	return p.root.Shutdown(ctx)
}

func (p *OurSpanProcessor) ForceFlush(ctx context.Context) error {
	log.Println("on ForceFlush")
	return p.root.ForceFlush(ctx)
}
