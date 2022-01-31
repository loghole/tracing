package tracing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/loghole/tracing/internal/logtracer"
	"github.com/loghole/tracing/internal/spanprocessor"
)

const _defaultTracerName = "github.com/loghole/tracing"

var _ trace.Tracer = new(Tracer)

// Tracer is the wrapper for `trace.Tracer` with span builder.
type Tracer struct {
	provider trace.TracerProvider
	tracer   trace.Tracer

	shutdown func(ctx context.Context) error
}

// NewTracer returns initialized Tracer with jaeger exporter.
//
// Example:
//
// func main() {
//     tracer, err := tracing.NewTracer(DefaultConfiguration(
//         "example-service",
//         "udp://127.0.0.1:6831",
//     ))
//     if err != nil {
//         log.Fatalf("init tracer: %v", err)
//     }
//
//     defer tracer.Close()
//
//     _, span := tracer.NewSpan().WithName("example-span").StartWithContext(context.Background())
//     defer span.End()
//
//     time.Sleep(time.Second)
// }
//
func NewTracer(configuration *Configuration) (*Tracer, error) {
	if err := configuration.validate(); err != nil {
		return nil, err
	}

	if configuration.Disabled {
		var (
			provider = logtracer.NewProvider()
			tracer   = provider.Tracer(_defaultTracerName)
		)

		otel.SetTracerProvider(provider)

		return &Tracer{provider: provider, tracer: tracer}, nil
	}

	endpoint, err := configuration.endpoint()
	if err != nil {
		return nil, err
	}

	exporter, err := jaeger.New(endpoint)
	if err != nil {
		return nil, fmt.Errorf("init jaeger exporter: %w", err)
	}

	processor := spanprocessor.NewSampled(
		tracesdk.NewBatchSpanProcessor(exporter, configuration.SpanProcessorOptions...),
		configuration.Sampler,
	)

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(processor),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			append(configuration.Attributes, semconv.ServiceNameKey.String(configuration.ServiceName))...,
		)),
	)

	otel.SetTracerProvider(provider)

	tracer := &Tracer{
		provider: provider,
		tracer:   provider.Tracer(_defaultTracerName),
		shutdown: provider.Shutdown,
	}

	return tracer, nil
}

// NewNoopTracer returns a Tracer that performs no operations.
// The Tracer and Spans created from the returned
// Tracer also perform no operations.
func NewNoopTracer() *Tracer {
	var (
		provider = trace.NewNoopTracerProvider()
		tracer   = provider.Tracer(_defaultTracerName)
	)

	return &Tracer{provider: provider, tracer: tracer}
}

// Start creates a span and a context.Context containing the newly-created span.
//
// If the context.Context provided in `ctx` contains a Span then the newly-created
// Span will be a child of that span, otherwise it will be a root span. This behavior
// can be overridden by providing `WithNewRoot()` as a SpanOption, causing the
// newly-created Span to be a root span even if `ctx` contains a Span.
//
// When creating a Span it is recommended to provide all known span attributes using
// the `WithAttributes()` SpanOption as samplers will only have access to the
// attributes provided when a Span is created.
//
// Any Span that is created MUST also be ended. This is the responsibility of the user.
// Implementations of this API may leak memory or other resources if Spans are not ended.
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// Close shuts down the span processors in the order they were registered.
func (t *Tracer) Close() error {
	if t.shutdown == nil {
		return nil
	}

	const timeout = time.Second * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return t.shutdown(ctx)
}

// NewSpan creates a SpanBuilder.
func (t *Tracer) NewSpan() SpanBuilder {
	return SpanBuilder{tracer: t.tracer}
}

// SpanBuilder is a helper for creating a span through a chain of calls.
type SpanBuilder struct {
	name    string
	tracer  trace.Tracer
	carrier propagation.TextMapCarrier
	options []trace.SpanStartOption
}

// WithName sets custom span name.
func (b SpanBuilder) WithName(name string) SpanBuilder {
	b.name = name

	return b
}

// ExtractMap reads tracecontext from the `map[string]string` carrier and set remote SpanContext for new span.
func (b SpanBuilder) ExtractMap(carrier map[string]string) SpanBuilder {
	if carrier == nil {
		return b
	}

	b.carrier = propagation.MapCarrier(carrier)

	return b
}

// ExtractHeaders reads tracecontext from the `http.Header` carrier and set remote SpanContext for new span.
func (b SpanBuilder) ExtractHeaders(carrier http.Header) SpanBuilder {
	if carrier == nil {
		return b
	}

	b.carrier = propagation.HeaderCarrier(carrier)

	return b
}

// Start creates a span.
func (b SpanBuilder) Start(ctx context.Context) *Span {
	_, span := b.StartWithContext(ctx)

	return span
}

// StartWithContext creates a span and a context.Context containing the newly-created span.
func (b SpanBuilder) StartWithContext(ctx context.Context) (context.Context, *Span) {
	if b.name == "" {
		b.name = callerName()
	}

	if b.carrier != nil {
		ctx = new(propagation.TraceContext).Extract(ctx, b.carrier)
	}

	ctx, span := b.tracer.Start(ctx, b.name, b.options...)

	return ctx, &Span{span: span, tracer: b.tracer}
}
