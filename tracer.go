package tracing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/loghole/tracing/internal/logtracer"
	"github.com/loghole/tracing/internal/spanprocessor"
)

const _defaultTracerName = "github.com/loghole/tracing"

var ErrInvalidConfiguration = errors.New("invalid configuration")

type Configuration struct { // nolint:govet // not need.
	ServiceName string
	Addr        string
	Disabled    bool
	BatchSize   int

	Sampler              tracesdk.Sampler
	Attributes           attribute.KeyValue
	SpanProcessorOptions []tracesdk.BatchSpanProcessorOption
}

type Tracer struct {
	provider trace.TracerProvider
	tracer   trace.Tracer

	shutdown func(ctx context.Context) error
}

func DefaultConfiguration(service, addr string) *Configuration {
	configuration := &Configuration{
		ServiceName: service,
		Addr:        addr,
		Disabled:    addr == "",
		Sampler:     tracesdk.AlwaysSample(),
	}

	return configuration
}

func NewTracer(configuration *Configuration) (*Tracer, error) {
	if configuration.Disabled {
		var (
			provider = logtracer.NewProvider()
			tracer   = provider.Tracer(_defaultTracerName)
		)

		otel.SetTracerProvider(provider)

		return &Tracer{provider: provider, tracer: tracer}, nil
	}

	var endpoint jaeger.EndpointOption

	u, err := url.Parse(configuration.Addr)
	if err != nil {
		return nil, fmt.Errorf("parse configuration addr: %w", err)
	}

	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		endpoint = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(u.String()))
	case "udp":
		endpoint = jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(u.Host),
			jaeger.WithAgentPort(u.Port()),
		)
	default:
		return nil, fmt.Errorf("%w: unknown addr scheme, supported [http, https, udp]", ErrInvalidConfiguration)
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
			semconv.ServiceNameKey.String(configuration.ServiceName),
		)),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			configuration.Attributes,
		)),
	)

	otel.SetTracerProvider(provider)

	tracer := &Tracer{
		provider: provider,
		tracer:   provider.Tracer(_defaultTracerName),
	}

	return tracer, nil
}

func (c *Tracer) Close() error {
	if c.shutdown == nil {
		return nil
	}

	const timeout = time.Second * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.shutdown(ctx)
}

func (c *Tracer) NewSpan() SpanBuilder {
	return SpanBuilder{tracer: c.tracer}
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

func (b SpanBuilder) ExtractBinary(carrier []byte) SpanBuilder {
	if carrier == nil {
		return b
	}

	// TODO: implement me.

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
