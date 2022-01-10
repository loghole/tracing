package tracing

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	baseSkipCallers = 3
)

type Span struct {
	tracer trace.Tracer
	span   trace.Span
}

func ChildSpan(ctx *context.Context) (s *Span) { // nolint:gocritic
	s = &Span{}

	if span := trace.SpanFromContext(*ctx); span != nil {
		*ctx, s.span = span.TracerProvider().Tracer(_defaultTracerName).Start(*ctx, callerName())
	}

	return s
}

// Finish is alias of End function.
func (s *Span) Finish() {
	s.End()
}

// End completes the Span. The Span is considered complete and ready to be
// delivered through the rest of the telemetry pipeline after this method
// is called. Therefore, updates to the Span are not allowed after this
// method has been called.
func (s *Span) End(options ...trace.SpanEndOption) {
	if s.span == nil && !s.span.IsRecording() {
		return
	}

	s.span.End(options...)
}

// AddEvent adds an event with the provided name and options.
func (s *Span) AddEvent(name string, options ...trace.EventOption) {
	if s.span == nil {
		return
	}

	s.span.AddEvent(name, options...)
}

// IsRecording returns the recording state of the Span. It will return
// true if the Span is active and events can be recorded.
func (s *Span) IsRecording() bool {
	if s.span == nil {
		return false
	}

	return s.span.IsRecording()
}

// RecordError will record err as an exception span event for this span. An
// additional call to SetStatus is required if the Status of the Span should
// be set to Error, as this method does not change the Span status. If this
// span is not being recorded or err is nil then this method does nothing.
func (s *Span) RecordError(err error, options ...trace.EventOption) {
	if s.span == nil {
		return
	}

	s.span.RecordError(err, options...)
}

// SpanContext returns the SpanContext of the Span. The returned SpanContext
// is usable even after the End method has been called for the Span.
func (s *Span) SpanContext() trace.SpanContext {
	if s.span == nil {
		return trace.SpanContext{}
	}

	return s.span.SpanContext()
}

// SetStatus sets the status of the Span in the form of a code and a
// description, overriding previous values set. The description is only
// included in a status when the code is for an error.
func (s *Span) SetStatus(code codes.Code, description string) {
	if s.span == nil {
		return
	}

	s.span.SetStatus(code, description)
}

// SetName sets the Span name.
func (s *Span) SetName(name string) {
	if s.span == nil {
		return
	}

	s.span.SetName(name)
}

// SetAttributes sets kv as attributes of the Span. If a key from kv
// already exists for an attribute of the Span it will be overwritten with
// the value contained in kv.
func (s *Span) SetAttributes(kv ...attribute.KeyValue) {
	if s.span == nil {
		return
	}

	s.span.SetAttributes(kv...)
}

// TracerProvider returns a TracerProvider that can be used to generate
// additional Spans on the same telemetry pipeline as the current Span.
func (s *Span) TracerProvider() trace.TracerProvider {
	if s.span == nil {
		return nil
	}

	return s.span.TracerProvider()
}

func (s *Span) SetTag(key string, value interface{}) *Span {
	if s.span == nil {
		return nil
	}

	switch val := value.(type) {
	case bool:
		s.span.SetAttributes(attribute.Bool(key, val))
	case []bool:
		s.span.SetAttributes(attribute.BoolSlice(key, val))
	case int:
		s.span.SetAttributes(attribute.Int(key, val))
	case []int:
		s.span.SetAttributes(attribute.IntSlice(key, val))
	case int64:
		s.span.SetAttributes(attribute.Int64(key, val))
	case []int64:
		s.span.SetAttributes(attribute.Int64Slice(key, val))
	case float64:
		s.span.SetAttributes(attribute.Float64(key, val))
	case []float64:
		s.span.SetAttributes(attribute.Float64Slice(key, val))
	case string:
		s.span.SetAttributes(attribute.String(key, val))
	case []string:
		s.span.SetAttributes(attribute.StringSlice(key, val))
	case fmt.Stringer:
		s.span.SetAttributes(attribute.Stringer(key, val))
	default:
		s.span.SetAttributes(attribute.String(key, fmt.Sprint(val)))
	}

	return s
}

func callerName() string {
	var pc [1]uintptr

	runtime.Callers(baseSkipCallers, pc[:])

	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return "unknown"
	}

	list := strings.Split(f.Name(), "/")

	return list[len(list)-1]
}
