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

type Tracer struct {
	provider trace.TracerProvider
	tracer   trace.Tracer

	shutdown func(ctx context.Context) error
}

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
	}

	return tracer, nil
}

func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

func (t *Tracer) Close() error {
	if t.shutdown == nil {
		return nil
	}

	const timeout = time.Second * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return t.shutdown(ctx)
}

func (t *Tracer) NewSpan() SpanBuilder {
	return SpanBuilder{tracer: t.tracer}
}

type SpanBuilder struct {
	name    string
	tracer  trace.Tracer
	carrier propagation.TextMapCarrier
	options []trace.SpanStartOption
}

func (b SpanBuilder) WithName(name string) SpanBuilder {
	b.name = name

	return b
}

func (b SpanBuilder) ExtractMap(carrier map[string]string) SpanBuilder {
	if carrier == nil {
		return b
	}

	b.carrier = propagation.MapCarrier(carrier)

	return b
}

func (b SpanBuilder) ExtractHeaders(carrier http.Header) SpanBuilder {
	if carrier == nil {
		return b
	}

	b.carrier = propagation.HeaderCarrier(carrier)

	return b
}

func (b SpanBuilder) Start(ctx context.Context) *Span {
	_, span := b.StartWithContext(ctx)

	return span
}

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
