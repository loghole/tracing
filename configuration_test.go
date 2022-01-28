package tracing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestDefaultConfiguration(t *testing.T) {
	type args struct {
		service string
		addr    string
	}
	tests := []struct {
		name string
		args args
		want *Configuration
	}{
		{
			name: "enabled",
			args: args{
				service: "test",
				addr:    "udp://127.0.0.1:6543",
			},
			want: &Configuration{
				ServiceName: "test",
				Addr:        "udp://127.0.0.1:6543",
				Disabled:    false,
				Sampler:     tracesdk.AlwaysSample(),
			},
		},
		{
			name: "disabled",
			args: args{
				service: "test",
				addr:    "",
			},
			want: &Configuration{
				ServiceName: "test",
				Addr:        "",
				Disabled:    true,
				Sampler:     tracesdk.AlwaysSample(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DefaultConfiguration(tt.args.service, tt.args.addr), "DefaultConfiguration(%v, %v)", tt.args.service, tt.args.addr)
		})
	}
}

func TestConfiguration_validate(t *testing.T) {
	type fields struct {
		ServiceName string
		Addr        string
		Disabled    bool
		Sampler     tracesdk.Sampler
		Attributes  []attribute.KeyValue
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "pass",
			fields: fields{
				ServiceName: "test",
				Sampler:     tracesdk.AlwaysSample(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "empty service name",
			fields: fields{
				ServiceName: "",
				Sampler:     tracesdk.AlwaysSample(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid configuration: empty service name")
			},
		},
		{
			name: "empty sampler",
			fields: fields{
				ServiceName: "test",
				Sampler:     nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid configuration: sampler cannot be empty")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Configuration{
				ServiceName: tt.fields.ServiceName,
				Addr:        tt.fields.Addr,
				Disabled:    tt.fields.Disabled,
				Sampler:     tt.fields.Sampler,
				Attributes:  tt.fields.Attributes,
			}
			tt.wantErr(t, c.validate(), fmt.Sprintf("validate()"))
		})
	}
}

func TestConfiguration_endpoint(t *testing.T) {
	tests := []struct {
		name    string
		c       *Configuration
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "http",
			c:       &Configuration{Addr: "http://127.0.0.1:8080/trace"},
			wantErr: assert.NoError,
		},
		{
			name:    "udp",
			c:       &Configuration{Addr: "udp://127.0.0.1:6543"},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			c:    &Configuration{Addr: ""},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorIs(t, err, ErrInvalidConfiguration)

				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.endpoint()
			if !tt.wantErr(t, err, fmt.Sprintf("endpoint()")) {
				return
			}

			assert.NotNil(t, got, "endpoint() is nil")
		})
	}
}
