package logtracer

import (
	"context"
	"encoding/binary"
	"math/rand"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type LogTraceProvider struct{}

func NewProvider() trace.TracerProvider {
	return &LogTraceProvider{}
}

func (p *LogTraceProvider) Tracer(_ string, _ ...trace.TracerOption) trace.Tracer {
	logtracer := &LogTracer{}

	seedGenerator := NewRand(time.Now().UnixNano())

	pool := sync.Pool{
		New: func() interface{} {
			return rand.NewSource(seedGenerator.Int63())
		},
	}

	logtracer.randomNumber = func() uint64 {
		var (
			generator = pool.Get().(rand.Source)
			number    = uint64(generator.Int63())
		)

		pool.Put(generator)

		return number
	}

	return logtracer
}

type LogTracer struct {
	randomNumber func() uint64
}

func (t *LogTracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	var result *LogSpan

	if span, ok := interface{}(trace.SpanFromContext(ctx)).(*LogSpan); ok {
		result = &LogSpan{traceID: span.traceID, spanID: t.generateSpanID(), tracer: t}
	} else {
		result = &LogSpan{traceID: t.generateTraceID(), spanID: t.generateSpanID(), tracer: t}
	}

	return trace.ContextWithSpan(ctx, result), result
}

func (t *LogTracer) generateTraceID() trace.TraceID {
	var id trace.TraceID

	binary.LittleEndian.PutUint64(id[:8], t.randomNumber())
	binary.LittleEndian.PutUint64(id[8:], t.randomNumber())

	return id
}

func (t *LogTracer) generateSpanID() trace.SpanID {
	var id trace.SpanID

	binary.LittleEndian.PutUint64(id[:], t.randomNumber())

	return id
}

type LogSpan struct {
	tracer  *LogTracer
	traceID trace.TraceID
	spanID  trace.SpanID
}

func (s *LogSpan) End(_ ...trace.SpanEndOption)                {}
func (s *LogSpan) AddEvent(_ string, _ ...trace.EventOption)   {}
func (s *LogSpan) IsRecording() bool                           { return false }
func (s *LogSpan) RecordError(_ error, _ ...trace.EventOption) {}
func (s *LogSpan) SetStatus(_ codes.Code, _ string)            {}
func (s *LogSpan) SetName(_ string)                            {}
func (s *LogSpan) SetAttributes(_ ...attribute.KeyValue)       {}
func (s *LogSpan) TracerProvider() trace.TracerProvider        { return trace.NewNoopTracerProvider() }

func (s *LogSpan) SpanContext() trace.SpanContext {
	return trace.SpanContext{}.
		WithTraceID(s.traceID).
		WithSpanID(s.spanID)
}
