package tracing

import (
	_ "go.opentelemetry.io/otel/sdk/trace"
)

//go:generate mockgen --build_flags=--mod=mod -destination mocks/otel.go -package mocks go.opentelemetry.io/otel/sdk/trace SpanProcessor
