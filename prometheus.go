package tracing

import (
	"github.com/loghole/tracing/internal/metrics"
)

func EnablePrometheusMetrics() error {
	return metrics.Register()
}
