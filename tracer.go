package tracing

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/gadavy/tracing/internal"
)

type Tracer struct {
	opentracing.Tracer
	closer io.Closer
}

func DefaultConfiguration(service, url string) *config.Configuration {
	return &config.Configuration{
		ServiceName: service,
		Disabled:    url == "",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: time.Second,
			LocalAgentHostPort:  url,
		},
	}
}

func NewTracer(configuration *config.Configuration, options ...config.Option) (*Tracer, error) {
	tracer, closer, err := configuration.NewTracer(options...)
	if err != nil {
		return nil, err
	}

	if _, ok := tracer.(*opentracing.NoopTracer); ok {
		return &Tracer{Tracer: internal.NewLogTracer(), closer: closer}, nil
	}

	return &Tracer{Tracer: tracer, closer: closer}, nil
}

func (c *Tracer) OpenTracer() opentracing.Tracer {
	return c.Tracer
}

func (c *Tracer) Close() error {
	if c.closer == nil {
		return nil
	}

	return c.closer.Close()
}

func (c *Tracer) NewSpan() SpanBuilder {
	return SpanBuilder{tracer: c.Tracer}
}

type SpanBuilder struct {
	name    string
	tracer  opentracing.Tracer
	options []opentracing.StartSpanOption
}

func (b SpanBuilder) WithName(name string) SpanBuilder {
	b.name = name

	return b
}

func (b SpanBuilder) ExtractMap(carrier map[string]string) SpanBuilder {
	if carrier == nil {
		return b
	}

	spanCtx, err := b.tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(carrier))
	if err == nil {
		b.options = append(b.options, opentracing.FollowsFrom(spanCtx))
	}

	return b
}

func (b SpanBuilder) ExtractBinary(carrier []byte) SpanBuilder {
	if carrier == nil {
		return b
	}

	spanCtx, err := b.tracer.Extract(opentracing.Binary, bytes.NewReader(carrier))
	if err == nil {
		b.options = append(b.options, opentracing.FollowsFrom(spanCtx))
	}

	return b
}

func (b SpanBuilder) ExtractHeaders(carrier http.Header) SpanBuilder {
	if carrier == nil {
		return b
	}

	spanCtx, err := b.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(carrier))
	if err == nil {
		b.options = append(b.options, opentracing.FollowsFrom(spanCtx))
	}

	return b
}

func (b SpanBuilder) Build() *Span {
	if b.name == "" {
		b.name = callerName()
	}

	return &Span{span: b.tracer.StartSpan(b.name, b.options...)}
}

func (b SpanBuilder) BuildWithContext(ctx context.Context) (*Span, context.Context) {
	span := b.Build()
	ctx = span.Context(ctx)

	return span, ctx
}
