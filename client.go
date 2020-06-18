package tracing

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type Config struct {
	// Адрес jaeger-agent. Example localhost:6831
	URI         string
	Enabled     bool
	ServiceName string
}

type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func NewTracer(cfg *Config) (*Tracer, error) {
	configuration := config.Configuration{
		ServiceName: cfg.ServiceName,
		Disabled:    !cfg.Enabled,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: time.Second,
			LocalAgentHostPort:  cfg.URI,
		},
	}

	tracer, closer, err := configuration.NewTracer(config.PoolSpans(true))
	if err != nil {
		return nil, err
	}

	return &Tracer{tracer: tracer, closer: closer}, nil
}

func (c *Tracer) OpenTracer() opentracing.Tracer {
	return c.tracer
}

func (c *Tracer) Close() error {
	return c.closer.Close()
}

func (c *Tracer) NewSpan() SpanBuilder {
	return SpanBuilder{tracer: c.tracer}
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

func (b SpanBuilder) Build() Span {
	return Span{span: b.tracer.StartSpan(b.name, b.options...)}
}
