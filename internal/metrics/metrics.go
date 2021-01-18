package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// nolint:gochecknoglobals // using metrics in clients.
var (
	logsCounters = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "logs_total",
		Help:        "Number of logs",
		ConstLabels: nil,
	}, []string{"level"})

	DebugLogsCounter = logsCounters.WithLabelValues("debug")
	InfoLogsCounter  = logsCounters.WithLabelValues("info")
	WarnLogsCounter  = logsCounters.WithLabelValues("warn")
	ErrorLogsCounter = logsCounters.WithLabelValues("error")

	outputRequestsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "output_requests_total",
		Help:        "Number of output requests",
		ConstLabels: nil,
	}, []string{"type", "status"})

	GRPCSuccessOutputReqCounter = outputRequestsCounter.WithLabelValues("grpc", "success")
	GRPCFailedOutputReqCounter  = outputRequestsCounter.WithLabelValues("grpc", "failed")

	HTTPSuccessOutputReqCounter = outputRequestsCounter.WithLabelValues("http", "success")
	HTTPFailedOutputReqCounter  = outputRequestsCounter.WithLabelValues("http", "failed")

	STANSuccessOutputReqCounter = outputRequestsCounter.WithLabelValues("stan", "success")
	STANFailedOutputReqCounter  = outputRequestsCounter.WithLabelValues("stan", "failed")

	inputRequestsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "input_requests_total",
		Help:        "Number of input requests",
		ConstLabels: nil,
	}, []string{"type", "status"})

	GRPCSuccessInputReqCounter = inputRequestsCounter.WithLabelValues("grpc", "success")
	GRPCFailedInputReqCounter  = inputRequestsCounter.WithLabelValues("grpc", "failed")

	HTTPSuccessInputReqCounter = inputRequestsCounter.WithLabelValues("http", "success")
	HTTPFailedInputReqCounter  = inputRequestsCounter.WithLabelValues("http", "failed")

	StanSuccessInputReqCounter = inputRequestsCounter.WithLabelValues("stan", "success")
	StanFailedInputReqCounter  = inputRequestsCounter.WithLabelValues("stan", "failed")
)

func Register() error {
	if err := prometheus.Register(logsCounters); err != nil {
		return fmt.Errorf("register logs counter: %w", err)
	}

	if err := prometheus.Register(outputRequestsCounter); err != nil {
		return fmt.Errorf("register output requests counter: %w", err)
	}

	if err := prometheus.Register(inputRequestsCounter); err != nil {
		return fmt.Errorf("register input requests counter: %w", err)
	}

	return nil
}
