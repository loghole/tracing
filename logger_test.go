package tracing

import (
	"context"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestTrace_ContextWithAction(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "#1",
			ctx:      ContextWithAction(context.Background(), "some action"),
			expected: "some action",
		},
		{
			name:     "#2",
			ctx:      context.Background(),
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ActionFromContext(tt.ctx))
		})
	}
}

func TestTraceLogger_withErrorTag(t *testing.T) {
	tests := []struct {
		name    string
		span    opentracing.Span
		withTag bool
	}{
		{
			name:    "WithSpan",
			span:    mocktracer.New().StartSpan("test"),
			withTag: true,
		},
		{
			name:    "WithoutSpan",
			span:    nil,
			withTag: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := opentracing.ContextWithSpan(context.TODO(), tt.span)

			withErrorTag(ctx)

			var result bool

			if tt.withTag {
				mockSpan, ok := tt.span.(*mocktracer.MockSpan)
				if !ok {
					t.Error("expected mocktracer.MockSpan")
				}

				result, ok = mockSpan.Tag("error").(bool)
				if !ok {
					t.Error("expected bool tag value")
				}
			}

			assert.Equal(t, tt.withTag, result)
		})
	}
}
