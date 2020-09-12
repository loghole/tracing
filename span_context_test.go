package tracing

import (
	"context"
	"net/http"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestInjectMap(t *testing.T) {
	tests := []struct {
		name     string
		span     opentracing.Span
		expected map[string]string
	}{
		{
			name:     "InjectMap",
			span:     mocktracer.New().StartSpan("test"),
			expected: map[string]string{"mockpfx-ids-sampled": "true", "mockpfx-ids-spanid": "44", "mockpfx-ids-traceid": "43"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := opentracing.ContextWithSpan(context.TODO(), tt.span)

			carrier := make(map[string]string)

			if err := InjectMap(ctx, carrier); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, carrier)
		})
	}
}

func TestInjectHeaders(t *testing.T) {
	mocktracer.New().Reset()

	tests := []struct {
		name     string
		span     opentracing.Span
		expected http.Header
	}{
		{
			name:     "InjectHeaders",
			span:     mocktracer.New().StartSpan("test"),
			expected: http.Header{"Mockpfx-Ids-Sampled": []string{"true"}, "Mockpfx-Ids-Spanid": []string{"46"}, "Mockpfx-Ids-Traceid": []string{"45"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := opentracing.ContextWithSpan(context.TODO(), tt.span)

			carrier := http.Header{}

			if err := InjectHeaders(ctx, carrier); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, carrier)
		})
	}
}
