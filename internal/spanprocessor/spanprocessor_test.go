package spanprocessor

import (
	"context"
	"encoding/binary"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/loghole/tracing/mocks"
)

func TestSampled_OnStart(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type args struct {
		makeProcessor func() tracesdk.SpanProcessor
		makeSpan      func() tracesdk.ReadWriteSpan
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "pass",
			args: args{
				makeProcessor: func() tracesdk.SpanProcessor {
					processor := mocks.NewMockSpanProcessor(ctrl)
					processor.EXPECT().OnStart(ctx, mocks.NewMockReadWriteSpan(ctrl)).Times(1)

					return processor
				},
				makeSpan: func() tracesdk.ReadWriteSpan {
					span := mocks.NewMockReadWriteSpan(ctrl)
					span.EXPECT().SpanContext().Times(1).Return(trace.SpanContext{})

					return span
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			processor := NewSampled(tt.args.makeProcessor(), tracesdk.AlwaysSample())

			processor.OnStart(ctx, tt.args.makeSpan())
		})
	}
}

func TestSampled_OnEnd(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type args struct {
		makeProcessor func() tracesdk.SpanProcessor
		makeSpan      func() tracesdk.ReadWriteSpan
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "pass",
			args: args{
				makeProcessor: func() tracesdk.SpanProcessor {
					processor := mocks.NewMockSpanProcessor(ctrl)
					processor.EXPECT().OnStart(ctx, mocks.NewMockReadWriteSpan(ctrl))

					processor.EXPECT().OnEnd(gomock.Any())

					return processor
				},
				makeSpan: func() tracesdk.ReadWriteSpan {
					span := mocks.NewMockReadWriteSpan(ctrl)
					span.EXPECT().SpanContext().Return(trace.SpanContext{}).Times(4)
					span.EXPECT().Attributes().Return([]attribute.KeyValue{}).Times(2)
					span.EXPECT().Status().Return(tracesdk.Status{})
					span.EXPECT().EndTime().Return(time.Now())
					span.EXPECT().Name().Return("")
					span.EXPECT().SpanKind().Return(trace.SpanKind(0))
					span.EXPECT().Links().Return([]tracesdk.Link{})

					return span
				},
			},
		},
		{
			name: "pass ended",
			args: args{
				makeProcessor: func() tracesdk.SpanProcessor {
					processor := mocks.NewMockSpanProcessor(ctrl)
					processor.EXPECT().OnStart(ctx, mocks.NewMockReadWriteSpan(ctrl))
					processor.EXPECT().OnEnd(gomock.Any())

					return processor
				},
				makeSpan: func() tracesdk.ReadWriteSpan {
					span := mocks.NewMockReadWriteSpan(ctrl)
					span.EXPECT().SpanContext().Return(trace.SpanContext{}).Times(4)
					span.EXPECT().EndTime().Return(time.Now())
					span.EXPECT().Attributes().Return([]attribute.KeyValue{}).Times(2)
					span.EXPECT().Status().Return(tracesdk.Status{})
					span.EXPECT().Name().Return("")
					span.EXPECT().SpanKind().Return(trace.SpanKind(0))
					span.EXPECT().Links().Return([]tracesdk.Link{})

					return span
				},
			},
		},
		{
			name: "pass ended with status err",
			args: args{
				makeProcessor: func() tracesdk.SpanProcessor {
					processor := mocks.NewMockSpanProcessor(ctrl)
					processor.EXPECT().OnStart(ctx, mocks.NewMockReadWriteSpan(ctrl))
					processor.EXPECT().OnEnd(gomock.Any())

					return processor
				},
				makeSpan: func() tracesdk.ReadWriteSpan {
					span := mocks.NewMockReadWriteSpan(ctrl)
					span.EXPECT().SpanContext().Return(trace.SpanContext{}).Times(3)
					span.EXPECT().EndTime().Return(time.Now())
					span.EXPECT().Status().Return(tracesdk.Status{Code: codes.Error})

					return span
				},
			},
		},
		{
			name: "pass ended with attr err",
			args: args{
				makeProcessor: func() tracesdk.SpanProcessor {
					processor := mocks.NewMockSpanProcessor(ctrl)
					processor.EXPECT().OnStart(ctx, mocks.NewMockReadWriteSpan(ctrl))
					processor.EXPECT().OnEnd(gomock.Any())

					return processor
				},
				makeSpan: func() tracesdk.ReadWriteSpan {
					span := mocks.NewMockReadWriteSpan(ctrl)
					span.EXPECT().SpanContext().Return(trace.SpanContext{}).Times(3)
					span.EXPECT().EndTime().Return(time.Now())
					span.EXPECT().Status().Return(tracesdk.Status{})
					span.EXPECT().Attributes().Return([]attribute.KeyValue{attribute.Bool("error", true)})

					return span
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			processor := NewSampled(tt.args.makeProcessor(), tracesdk.AlwaysSample())

			span := tt.args.makeSpan()

			processor.OnStart(ctx, span)
			processor.OnEnd(span)
		})
	}
}

func TestSampled_Shutdown(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	tests := []struct {
		name          string
		makeProcessor func() *Sampled
		wantErr       bool
	}{
		{
			name: "pass",
			makeProcessor: func() *Sampled {
				processor := mocks.NewMockSpanProcessor(ctrl)
				processor.EXPECT().Shutdown(ctx).Return(nil)
				processor.EXPECT().OnEnd(gomock.Any())

				sampled := NewSampled(processor, tracesdk.AlwaysSample())
				sampled.traces[trace.TraceID{}] = &traceWrapper{
					spans: map[trace.SpanID]spanWrapper{
						trace.SpanID{}: {
							span: &NoopSpan{},
							ctx:  ctx,
						},
					},
					isFinished: true,
					_hasError:  false,
				}

				return sampled
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.makeProcessor().Shutdown(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Sampled.Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSampled_ForceFlush(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	tests := []struct {
		name          string
		makeProcessor func() *Sampled
		wantErr       bool
	}{
		{
			name: "pass",
			makeProcessor: func() *Sampled {
				processor := mocks.NewMockSpanProcessor(ctrl)
				processor.EXPECT().ForceFlush(ctx).Return(nil)
				processor.EXPECT().OnEnd(gomock.Any())

				sampled := NewSampled(processor, tracesdk.AlwaysSample())
				sampled.traces[trace.TraceID{}] = &traceWrapper{
					spans: map[trace.SpanID]spanWrapper{
						trace.SpanID{}: {
							span: &NoopSpan{},
							ctx:  ctx,
						},
					},
					isFinished: true,
					_hasError:  false,
				}

				return sampled
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.makeProcessor().ForceFlush(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Sampled.ForceFlush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSampled_BaseWork(t *testing.T) {
	var (
		generator = NewSpanGenerator()
		processor = NewSampled(NoopSpanProcessor{}, tracesdk.AlwaysSample())
	)

	for i := 0; i < 1000; i++ {
		ctx, span := generator.NewSpanWithContext()
		processor.OnStart(ctx, span)
		processor.OnEnd(span)
	}

	assert.Equal(t, map[trace.TraceID]*traceWrapper{}, processor.traces)
}

func TestSampled_SendSpan(t *testing.T) {
	type args struct {
		sampler tracesdk.Sampler
	}

	tests := []struct {
		name    string
		args    args
		setter  func(ctx context.Context, span trace.Span)
		checker func(recorder *tracetest.SpanRecorder)
	}{
		{
			name: "AlwaysSample span without error",
			args: args{
				sampler: tracesdk.AlwaysSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 1, "recorder.Ended() != 1")
			},
		},
		{
			name: "AlwaysSample span with error status",
			args: args{
				sampler: tracesdk.AlwaysSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {
				span.SetStatus(codes.Error, "some error")
			},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 1, "recorder.Ended() != 1")
			},
		},
		{
			name: "AlwaysSample span with error attribute",
			args: args{
				sampler: tracesdk.AlwaysSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {
				span.SetAttributes(attribute.Bool("error", true))
			},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 1, "recorder.Ended() != 1")
			},
		},

		// NeverSample.
		{
			name: "NeverSample span without error",
			args: args{
				sampler: tracesdk.NeverSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 0, "recorder.Ended() != 0")
			},
		},
		{
			name: "NeverSample span with error status",
			args: args{
				sampler: tracesdk.NeverSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {
				span.SetStatus(codes.Error, "some error")
			},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 1, "recorder.Ended() != 1")
			},
		},
		{
			name: "NeverSample span with error attribute",
			args: args{
				sampler: tracesdk.NeverSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {
				span.SetAttributes(attribute.Bool("error", true))
			},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 1, "recorder.Ended() != 1")
			},
		},

		// Child span.
		{
			name: "ChildSpan finished",
			args: args{
				sampler: tracesdk.AlwaysSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {
				_, span2 := span.TracerProvider().Tracer("").Start(ctx, "next")
				span2.End()
			},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 2, "recorder.Ended() != 2")
			},
		},
		{
			name: "ChildSpan not finished",
			args: args{
				sampler: tracesdk.AlwaysSample(),
			},
			setter: func(ctx context.Context, span trace.Span) {
				_, _ = span.TracerProvider().Tracer("").Start(ctx, "next")
			},
			checker: func(recorder *tracetest.SpanRecorder) {
				assert.Len(t, recorder.Ended(), 1, "recorder.Ended() != 1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				recorder  = tracetest.NewSpanRecorder()
				processor = NewSampled(recorder, tt.args.sampler)
				tracer    = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor)).Tracer("")
			)

			ctx, span := tracer.Start(context.Background(), "test")
			tt.setter(ctx, span)
			span.End()

			tt.checker(recorder)
		})
	}
}

type NoopSpan struct {
	tracesdk.ReadWriteSpan
	context trace.SpanContext
}

func (s *NoopSpan) End(_ ...trace.SpanEndOption)                {}
func (s *NoopSpan) AddEvent(_ string, _ ...trace.EventOption)   {}
func (s *NoopSpan) IsRecording() bool                           { return true }
func (s *NoopSpan) RecordError(_ error, _ ...trace.EventOption) {}
func (s *NoopSpan) SetStatus(_ codes.Code, _ string)            {}
func (s *NoopSpan) SetName(_ string)                            {}
func (s *NoopSpan) SetAttributes(_ ...attribute.KeyValue)       {}
func (s *NoopSpan) TracerProvider() trace.TracerProvider        { return trace.NewNoopTracerProvider() }
func (s *NoopSpan) SpanContext() trace.SpanContext              { return s.context }
func (s *NoopSpan) Attributes() []attribute.KeyValue            { return []attribute.KeyValue{} }
func (s *NoopSpan) Parent() trace.SpanContext                   { return trace.SpanContext{} }
func (s *NoopSpan) SpanKind() trace.SpanKind                    { return 0 }
func (s *NoopSpan) StartTime() time.Time                        { return time.Time{} }
func (s *NoopSpan) EndTime() time.Time                          { return time.Now() }
func (s *NoopSpan) Links() []tracesdk.Link                      { return []tracesdk.Link{} }
func (s *NoopSpan) Events() []tracesdk.Event                    { return []tracesdk.Event{} }
func (s *NoopSpan) Status() tracesdk.Status                     { return tracesdk.Status{} }
func (s *NoopSpan) Name() string                                { return "" }

type SpanGenerator struct {
	traceID uint64
	spanID  uint64
}

func NewSpanGenerator() *SpanGenerator {
	return &SpanGenerator{}
}

func (s *SpanGenerator) NewSpan() *NoopSpan {
	var (
		traceID trace.TraceID
		spanID  trace.SpanID
	)

	binary.LittleEndian.PutUint64(traceID[:8], atomic.AddUint64(&s.traceID, 1))
	binary.LittleEndian.PutUint64(spanID[:8], atomic.AddUint64(&s.spanID, 1))

	sc := trace.SpanContext{}
	sc = sc.WithTraceID(traceID)
	sc = sc.WithSpanID(spanID)

	return &NoopSpan{context: sc}
}

func (s *SpanGenerator) NewSpanWithContext() (context.Context, *NoopSpan) {
	span := s.NewSpan()

	return trace.ContextWithSpan(context.Background(), span), span
}

func (s *SpanGenerator) Reset() {
	atomic.StoreUint64(&s.traceID, 0)
	atomic.StoreUint64(&s.spanID, 0)
}

type NoopSpanProcessor struct{}

func (n NoopSpanProcessor) OnStart(parent context.Context, s tracesdk.ReadWriteSpan) {}
func (n NoopSpanProcessor) OnEnd(s tracesdk.ReadOnlySpan)                            {}
func (n NoopSpanProcessor) Shutdown(ctx context.Context) error                       { return nil }
func (n NoopSpanProcessor) ForceFlush(ctx context.Context) error                     { return nil }

func BenchmarkSampled(b *testing.B) {

	b.Run("OnStart", func(b *testing.B) {
		var (
			generator = NewSpanGenerator()
			processor = NewSampled(NoopSpanProcessor{}, tracesdk.AlwaysSample())
		)

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ctx, span := generator.NewSpanWithContext()
				processor.OnStart(ctx, span)
			}
		})
	})

	b.Run("OnEnd", func(b *testing.B) {
		var (
			generator = NewSpanGenerator()
			processor = NewSampled(NoopSpanProcessor{}, tracesdk.AlwaysSample())
		)

		for i := 0; i < b.N; i++ {
			ctx, span := generator.NewSpanWithContext()
			processor.OnStart(ctx, span)
		}

		generator.Reset()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				span := generator.NewSpan()

				processor.OnEnd(span)
			}
		})
	})

	b.Run("OnStart|OnEnd", func(b *testing.B) {
		var (
			generator = NewSpanGenerator()
			processor = NewSampled(NoopSpanProcessor{}, tracesdk.AlwaysSample())
		)

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ctx, span := generator.NewSpanWithContext()
				go processor.OnEnd(span)

				processor.OnStart(ctx, span)
			}
		})
	})
}
