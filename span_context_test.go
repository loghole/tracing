package tracing

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"

	"github.com/loghole/tracing/mocks"
)

func TestInjectMap(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "empty context",
			args: args{
				ctx: context.Background(),
			},
			want: map[string]string{},
		},
		{
			name: "pass",
			args: args{
				ctx: mocks.NewContextWithMockSpan(context.Background(), 1, 2),
			},
			want: map[string]string{"traceparent": "00-01000000000000000000000000000000-0200000000000000-00"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carrier := map[string]string{}

			InjectMap(tt.args.ctx, carrier)

			assert.Equal(t, carrier, tt.want)
		})
	}
}

func TestInjectHeaders(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		{
			name: "empty context",
			args: args{
				ctx: context.Background(),
			},
			want: http.Header{},
		},
		{
			name: "pass",
			args: args{
				ctx: mocks.NewContextWithMockSpan(context.Background(), 1, 2),
			},
			want: http.Header{"Traceparent": []string{"00-01000000000000000000000000000000-0200000000000000-00"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carrier := http.Header{}

			InjectHeaders(tt.args.ctx, carrier)

			assert.Equal(t, carrier, tt.want)
		})
	}
}

func TestSpanContextFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want func() trace.SpanContext
	}{
		{
			name: "pass",
			args: args{
				ctx: mocks.NewContextWithMockSpan(context.Background(), 1, 2),
			},
			want: func() trace.SpanContext {
				sc := trace.SpanContext{}
				sc = sc.WithTraceID(trace.TraceID{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
				sc = sc.WithSpanID(trace.SpanID{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})

				return sc
			},
		},
		{
			name: "empty context",
			args: args{
				ctx: context.Background(),
			},
			want: func() trace.SpanContext {
				return trace.SpanContext{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want(), SpanContextFromContext(tt.args.ctx), "SpanContextFromContext(%v)", tt.args.ctx)
		})
	}
}

func TestSpanFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want trace.Span
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpanFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpanFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
