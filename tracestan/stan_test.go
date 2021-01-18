package tracestan

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func Test_reflectCb(t *testing.T) {
	tests := []struct {
		name        string
		args        Handler
		wantArgType reflect.Type
		wantErr     bool
		wantErrStr  string
	}{
		{
			name:        "pass",
			args:        func(ctx context.Context, status *errdetails.DebugInfo) error { return nil },
			wantArgType: reflect.TypeOf(new(errdetails.DebugInfo)),
			wantErr:     false,
		},
		{
			name:       "invalid callback type",
			args:       "callback",
			wantErr:    true,
			wantErrStr: "invalid callback args: handler has to be a func",
		},
		{
			name:       "invalid num of args",
			args:       func(status *errdetails.DebugInfo) error { return nil },
			wantErr:    true,
			wantErrStr: "invalid callback args: num of arguments should be 2",
		},
		{
			name:       "invalid first arg",
			args:       func(ctx string, status *errdetails.DebugInfo) error { return nil },
			wantErr:    true,
			wantErrStr: "invalid callback args: first arg must be context",
		},
		{
			name:       "invalid second arg",
			args:       func(ctx context.Context, status struct{ Name string }) error { return nil },
			wantErr:    true,
			wantErrStr: "invalid callback args: second arg must implement proto.Message",
		},
		{
			name:       "invalid return num arg",
			args:       func(ctx context.Context, status *errdetails.DebugInfo) (string, error) { return "", nil },
			wantErr:    true,
			wantErrStr: "invalid callback result: num of return args should be 1",
		},
		{
			name:       "invalid return arg type",
			args:       func(ctx context.Context, status *errdetails.DebugInfo) string { return "" },
			wantErr:    true,
			wantErrStr: "invalid callback result: return type should be 'error'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCallback, gotArgType, err := reflectCb(tt.args)
			if tt.wantErr {
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				wantCallback := reflect.ValueOf(tt.args)

				assert.Equal(t, wantCallback, gotCallback)
				assert.Equal(t, tt.wantArgType, gotArgType)
			}
		})
	}
}

func TestClient_handler(t *testing.T) {
	type args struct {
		topic string
		cb    Handler
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrStr string
	}{
		{
			name: "pass",
			args: args{
				topic: "test-topic",
				cb:    func(ctx context.Context, info *errdetails.DebugInfo) error { return nil },
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				topic: "test-topic",
				cb:    func(info *errdetails.DebugInfo) error { return nil },
			},
			wantErr:    true,
			wantErrStr: "invalid callback args: num of arguments should be 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}

			got, err := c.buildHandler(tt.args.topic, tt.args.cb)
			if tt.wantErr {
				assert.EqualError(t, err, tt.wantErrStr)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
