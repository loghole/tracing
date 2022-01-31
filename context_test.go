package tracing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_noCancelCtx_Deadline(t *testing.T) {
	type fields struct {
		parent context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
		want1  bool
	}{
		{
			name: "pass",
			fields: fields{
				parent: context.Background(),
			},
			want:  time.Time{},
			want1: false,
		},
		{
			name: "pass with context canceled",
			fields: fields{
				parent: canceledContext(),
			},
			want:  time.Time{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithoutCancel(tt.fields.parent)

			got, got1 := ctx.Deadline()
			assert.Equalf(t, tt.want, got, "Deadline()")
			assert.Equalf(t, tt.want1, got1, "Deadline()")
		})
	}
}

func Test_noCancelCtx_Done(t *testing.T) {
	type fields struct {
		parent context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   <-chan struct{}
	}{
		{
			name: "pass",
			fields: fields{
				parent: context.Background(),
			},
			want: nil,
		},
		{
			name: "pass with context canceled",
			fields: fields{
				parent: canceledContext(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithoutCancel(tt.fields.parent)

			assert.Equalf(t, tt.want, ctx.Done(), "Done()")
		})
	}
}

func Test_noCancelCtx_Err(t *testing.T) {
	type fields struct {
		parent context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "pass",
			fields: fields{
				parent: context.Background(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "pass with context canceled",
			fields: fields{
				parent: canceledContext(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithoutCancel(tt.fields.parent)

			tt.wantErr(t, ctx.Err(), fmt.Sprintf("Err()"))
		})
	}
}

func Test_noCancelCtx_Value(t *testing.T) {
	type fields struct {
		parent context.Context
	}
	type args struct {
		key interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "pass",
			fields: fields{
				parent: context.Background(),
			},
			args: args{
				key: "key",
			},
			want: nil,
		},
		{
			name: "pass with context canceled",
			fields: fields{
				parent: canceledContext(),
			},
			args: args{
				key: "key",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithoutCancel(tt.fields.parent)

			assert.Equalf(t, tt.want, ctx.Value(tt.args.key), "Value(%v)", tt.args.key)
		})
	}
}

func TestContextWithoutCancel(t *testing.T) {
	assert.Panics(t, func() {
		ContextWithoutCancel(nil)
	})
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return ctx
}
