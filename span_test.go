package tracing

import (
	"context"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestChildSpan(t *testing.T) {
	tests := []struct {
		name     string
		span     opentracing.Span
		wantSpan bool
	}{
		{
			name:     "WithSpan",
			span:     mocktracer.New().StartSpan("test"),
			wantSpan: true,
		},
		{
			name:     "WithoutSpan",
			span:     nil,
			wantSpan: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := opentracing.ContextWithSpan(context.TODO(), tt.span)
			tracer := ChildSpan(&ctx)

			if (tracer.span != nil) != tt.wantSpan {
				t.Errorf("tracer span exists = %v, expected = %v", tracer.span != nil, tt.wantSpan)
			}

			if tt.wantSpan {
				parent, ok := tt.span.(*mocktracer.MockSpan)
				if !ok {
					t.Error("expected mocktracer.MockSpan")
				}

				child, ok := tracer.span.(*mocktracer.MockSpan)
				if !ok {
					t.Error("expected mocktracer.MockSpan")
				}

				assert.Equal(t, parent.SpanContext.TraceID, child.SpanContext.TraceID)
				assert.Equal(t, parent.SpanContext.SpanID, child.ParentID)
			}
		})
	}
}

func TestFollowsSpan(t *testing.T) {
	tests := []struct {
		name     string
		span     opentracing.Span
		wantSpan bool
	}{
		{
			name:     "WithSpan",
			span:     mocktracer.New().StartSpan("test"),
			wantSpan: true,
		},
		{
			name:     "WithoutSpan",
			span:     nil,
			wantSpan: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := opentracing.ContextWithSpan(context.TODO(), tt.span)
			tracer := FollowsSpan(&ctx)

			if (tracer.span != nil) != tt.wantSpan {
				t.Errorf("tracer span exists = %v, expected = %v", tracer.span != nil, tt.wantSpan)
			}

			if tt.wantSpan {
				parent, ok := tt.span.(*mocktracer.MockSpan)
				if !ok {
					t.Error("expected mocktracer.MockSpan")
				}

				child, ok := tracer.span.(*mocktracer.MockSpan)
				if !ok {
					t.Error("expected mocktracer.MockSpan")
				}

				assert.Equal(t, parent.SpanContext.TraceID, child.SpanContext.TraceID)
				assert.Equal(t, parent.SpanContext.SpanID, child.ParentID)
			}
		})
	}
}

func TestSpan_SetTag(t *testing.T) {
	tests := []struct {
		name   string
		span   opentracing.Span
		tagKey string
		tagVal interface{}
	}{
		{
			name:   "#1",
			span:   mocktracer.New().StartSpan("test"),
			tagKey: "key",
			tagVal: "val",
		},
		{
			name:   "#2",
			span:   mocktracer.New().StartSpan("test"),
			tagKey: "key",
			tagVal: 123456567890,
		},
		{
			name:   "#3",
			span:   mocktracer.New().StartSpan("test"),
			tagKey: "key",
			tagVal: struct {
				Key string
				Val int
			}{
				Key: "1234567890",
				Val: 1234567890,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := opentracing.ContextWithSpan(context.TODO(), tt.span)
			span := ChildSpan(&ctx).SetTag(tt.tagKey, tt.tagVal)

			parent, ok := tt.span.(*mocktracer.MockSpan)
			if !ok {
				t.Fatal("expected mocktracer.MockSpan")
			}

			child, ok := span.(*mocktracer.MockSpan)
			if !ok {
				t.Fatal("expected mocktracer.MockSpan")
			}

			assert.Equal(t, parent.SpanContext.TraceID, child.SpanContext.TraceID)
			assert.Equal(t, parent.SpanContext.SpanID, child.ParentID)
			assert.Equal(t, child.Tag(tt.tagKey), tt.tagVal)
		})
	}
}

func TestSpan_Finish(t *testing.T) {
	tests := []struct {
		name string
		span opentracing.Span
	}{
		{
			name: "WithSpan",
			span: mocktracer.New().StartSpan("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := opentracing.ContextWithSpan(context.TODO(), tt.span)
			tracer := ChildSpan(&ctx)
			tracer.Finish()
			tracer.Finish()

			parent, ok := tt.span.(*mocktracer.MockSpan)
			if !ok {
				t.Fatal("expected mocktracer.MockSpan")
			}

			child, ok := tracer.span.(*mocktracer.MockSpan)
			if !ok {
				t.Fatal("expected mocktracer.MockSpan")
			}

			assert.Equal(t, parent.SpanContext.TraceID, child.SpanContext.TraceID)
			assert.Equal(t, parent.SpanContext.SpanID, child.ParentID)
			assert.NotEqual(t, parent.StartTime, child.StartTime)
		})
	}
}
