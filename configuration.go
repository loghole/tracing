package tracing

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var ErrInvalidConfiguration = errors.New("invalid configuration")

// Configuration configures Tracer.
type Configuration struct { // nolint:govet // not need.
	ServiceName string
	Addr        string
	Disabled    bool

	Sampler              tracesdk.Sampler
	Attributes           []attribute.KeyValue
	SpanProcessorOptions []tracesdk.BatchSpanProcessorOption
}

// DefaultConfiguration returns base configuration with default params.
func DefaultConfiguration(service, addr string) *Configuration {
	configuration := &Configuration{
		ServiceName: service,
		Addr:        addr,
		Disabled:    addr == "",
		Sampler:     tracesdk.AlwaysSample(),
	}

	return configuration
}

func (c *Configuration) validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("%w: empty service name", ErrInvalidConfiguration)
	}

	if c.Sampler == nil {
		return fmt.Errorf("%w: smpler cannot be empty", ErrInvalidConfiguration)
	}

	return nil
}

func (c *Configuration) endpoint() (jaeger.EndpointOption, error) {
	u, err := url.Parse(c.Addr)
	if err != nil {
		return nil, fmt.Errorf("parse addr: %w", err)
	}

	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		return jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(u.String())), nil
	case "udp":
		return jaeger.WithAgentEndpoint(jaeger.WithAgentHost(u.Hostname()), jaeger.WithAgentPort(u.Port())), nil
	default:
		return nil, fmt.Errorf("%w: unknown addr scheme, supported [http, https, udp]", ErrInvalidConfiguration)
	}
}
