package tracing

import (
	"context"
	"net/http"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestLogKV(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
	}
	tests := []struct {
		name     string
		args     args
		withSpan bool
	}{
		{
			name: "WithSpan",
			args: args{
				ctx:   opentracing.ContextWithSpan(context.Background(), mocktracer.New().StartSpan("TestLogFields")),
				key:   "key",
				value: "value",
			},
			withSpan: true,
		},
		{
			name: "WithoutSpan",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				value: "value",
			},
			withSpan: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogKV(tt.args.ctx, tt.args.key, tt.args.value)

			if tt.withSpan {
				span, ok := SpanFromContext(tt.args.ctx).(*mocktracer.MockSpan)
				if !ok {
					t.Error("expected mocktracer.MockSpan")
				}

				logs := span.Logs()

				if len(logs) == 0 {
					t.Error("expected one log")
				}

				if len(logs[0].Fields) == 0 {
					t.Error("expected one log field")
				}

				assert.Equal(t, tt.args.key, logs[0].Fields[0].Key)
				assert.Equal(t, tt.args.value, logs[0].Fields[0].ValueString)
			}
		})
	}
}
func TestInjectMap(t *testing.T) {
	tests := []struct {
		name     string
		span     opentracing.Span
		expected map[string]string
	}{
		{
			name:     "InjectMap",
			span:     mocktracer.New().StartSpan("test"),
			expected: map[string]string{"mockpfx-ids-sampled": "true", "mockpfx-ids-spanid": "46", "mockpfx-ids-traceid": "45"},
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
			expected: http.Header{"Mockpfx-Ids-Sampled": []string{"true"}, "Mockpfx-Ids-Spanid": []string{"48"}, "Mockpfx-Ids-Traceid": []string{"47"}},
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
