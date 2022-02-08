package tracing

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func noopSpan() *Span {
	tracer := trace.NewNoopTracerProvider().Tracer("")
	_, span := tracer.Start(context.Background(), "")

	return &Span{
		tracer: tracer,
		span:   span,
	}
}

func TestSpan_Finish(t *testing.T) {
	tests := []struct {
		name string
		s    *Span
	}{
		{
			name: "pass",
			s:    noopSpan(),
		},
		{
			name: "pass with nil",
			s:    &Span{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Finish()
		})
	}
}

func TestSpan_End(t *testing.T) {
	type args struct {
		options []trace.SpanEndOption
	}
	tests := []struct {
		name string
		s    *Span
		args args
	}{
		{
			name: "pass",
			s:    noopSpan(),
			args: args{
				options: []trace.SpanEndOption{},
			},
		},
		{
			name: "pass with nil",
			s:    &Span{},
			args: args{
				options: []trace.SpanEndOption{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.End(tt.args.options...)
		})
	}
}

func TestSpan_AddEvent(t *testing.T) {
	type args struct {
		name    string
		options []trace.EventOption
	}
	tests := []struct {
		name string
		s    *Span
		args args
	}{
		{
			name: "pass",
			s:    noopSpan(),
			args: args{
				name: "test",
			},
		},
		{
			name: "pass with nil",
			s:    &Span{},
			args: args{
				name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.AddEvent(tt.args.name, tt.args.options...)
		})
	}
}

func TestSpan_IsRecording(t *testing.T) {
	tests := []struct {
		name string
		s    *Span
		want bool
	}{
		{
			name: "pass",
			s:    noopSpan(),
			want: false,
		},
		{
			name: "pass with nil",
			s:    &Span{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsRecording(); got != tt.want {
				t.Errorf("Span.IsRecording() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpan_RecordError(t *testing.T) {
	type args struct {
		err     error
		options []trace.EventOption
	}
	tests := []struct {
		name string
		s    *Span
		args args
	}{
		{
			name: "pass",
			s:    noopSpan(),
			args: args{
				err: fmt.Errorf("test"),
			},
		},
		{
			name: "pass with nil",
			s:    &Span{},
			args: args{
				err: fmt.Errorf("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.RecordError(tt.args.err, tt.args.options...)
		})
	}
}

func TestSpan_SpanContext(t *testing.T) {
	tests := []struct {
		name string
		s    *Span
		want trace.SpanContext
	}{
		{
			name: "pass with nil",
			s:    &Span{},
			want: trace.SpanContext{},
		},
		{
			name: "pass",
			s:    noopSpan(),
			want: trace.SpanContext{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.SpanContext()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSpan_SetStatus(t *testing.T) {
	type args struct {
		code        codes.Code
		description string
	}
	tests := []struct {
		name string
		s    *Span
		args args
	}{
		{
			name: "pass with nil",
			s:    &Span{},
			args: args{
				code:        12,
				description: "test",
			},
		},
		{
			name: "pass",
			s:    noopSpan(),
			args: args{
				code:        12,
				description: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetStatus(tt.args.code, tt.args.description)
		})
	}
}

func TestSpan_SetName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		s    *Span
		args args
	}{
		{
			name: "pass with nil",
			s:    &Span{},
			args: args{
				name: "test",
			},
		},
		{
			name: "pass",
			s:    noopSpan(),
			args: args{
				name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetName(tt.args.name)
		})
	}
}

func TestSpan_SetAttributes(t *testing.T) {
	type args struct {
		kv []attribute.KeyValue
	}
	tests := []struct {
		name string
		s    *Span
		args args
	}{
		{
			name: "pass with nil",
			s:    &Span{},
			args: args{
				kv: []attribute.KeyValue{
					{
						Key:   attribute.Key("test"),
						Value: attribute.Int64Value(12),
					},
				},
			},
		},
		{
			name: "pass",
			s:    noopSpan(),
			args: args{
				kv: []attribute.KeyValue{
					attribute.Int64Slice("test", []int64{12, 12}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetAttributes(tt.args.kv...)
		})
	}
}

func TestSpan_TracerProvider(t *testing.T) {
	tests := []struct {
		name string
		s    *Span
		want trace.TracerProvider
	}{
		{
			name: "pass with nil",
			s:    &Span{},
			want: trace.NewNoopTracerProvider(),
		},
		{
			name: "pass",
			s:    noopSpan(),
			want: trace.NewNoopTracerProvider(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.TracerProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Span.TracerProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpan_SetTag(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		s    *Span
		args args
	}{
		{
			name: "pass with nil",
			s:    &Span{},
			args: args{
				key:   "test",
				value: 12,
			},
		},
		{
			name: "pass",
			s:    noopSpan(),
			args: args{
				key:   "test",
				value: 11,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.SetTag(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.s) {
				t.Errorf("Span.SetTag() = %v, want %v", got, tt.s)
			}
		})
	}
}

func Test_attributeFromInterface(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want attribute.KeyValue
	}{
		{
			name: "Bool",
			args: args{
				key:   "key",
				value: true,
			},
			want: attribute.Bool("key", true),
		},
		{
			name: "BoolSlice",
			args: args{
				key:   "key",
				value: []bool{true, false},
			},
			want: attribute.BoolSlice("key", []bool{true, false}),
		},
		{
			name: "Int",
			args: args{
				key:   "key",
				value: 1,
			},
			want: attribute.Int("key", 1),
		},
		{
			name: "IntSlice",
			args: args{
				key:   "key",
				value: []int{1, 2, 3},
			},
			want: attribute.IntSlice("key", []int{1, 2, 3}),
		},
		{
			name: "Int64",
			args: args{
				key:   "key",
				value: int64(1),
			},
			want: attribute.Int64("key", 1),
		},
		{
			name: "Int64Slice",
			args: args{
				key:   "key",
				value: []int64{1, 2, 3},
			},
			want: attribute.Int64Slice("key", []int64{1, 2, 3}),
		},
		{
			name: "Float64",
			args: args{
				key:   "key",
				value: 1.1,
			},
			want: attribute.Float64("key", 1.1),
		},
		{
			name: "Float64Slice",
			args: args{
				key:   "key",
				value: []float64{1.1, 1.2},
			},
			want: attribute.Float64Slice("key", []float64{1.1, 1.2}),
		},
		{
			name: "String",
			args: args{
				key:   "key",
				value: "val",
			},
			want: attribute.String("key", "val"),
		},
		{
			name: "StringSlice",
			args: args{
				key:   "key",
				value: []string{"val", "val2"},
			},
			want: attribute.StringSlice("key", []string{"val", "val2"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, attributeFromInterface(tt.args.key, tt.args.value), "attributeFromInterface(%v, %v)", tt.args.key, tt.args.value)
		})
	}
}

func TestChildSpan(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		wantS *Span
	}{
		{
			name:  "pass",
			wantS: &Span{span: noopSpan().span},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := ChildSpan(&ctx); !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("ChildSpan() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func Test_callerName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "pass",
			want: "testing.tRunner",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := callerName(); got != tt.want {
				t.Errorf("callerName() = %v, want %v", got, tt.want)
			}
		})
	}
}
