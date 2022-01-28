package tracelog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/loghole/tracing/mocks"
)

func TestTraceLogger_Debug(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  context.Background(),
				args: []interface{}{"1", "2", "3"},
			},
			expected: "debug\t123\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				args: []interface{}{"some string", 1234567890},
			},
			expected: "debug\tsome string1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Debug(tt.args.ctx, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_Debugf(t *testing.T) {
	type args struct {
		ctx      context.Context
		template string
		args     []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      context.Background(),
				template: "some value: %s",
				args:     []interface{}{"value"},
			},
			expected: "debug\tsome value: value\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				template: "some int: %d",
				args:     []interface{}{1234567890},
			},
			expected: "debug\tsome int: 1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Debugf(tt.args.ctx, tt.args.template, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_Info(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  context.Background(),
				args: []interface{}{"1", "2", "3"},
			},
			expected: "info\t123\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				args: []interface{}{"some string", 1234567890},
			},
			expected: "info\tsome string1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Info(tt.args.ctx, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_Infof(t *testing.T) {
	type args struct {
		ctx      context.Context
		template string
		args     []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      context.Background(),
				template: "some value: %s",
				args:     []interface{}{"value"},
			},
			expected: "info\tsome value: value\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				template: "some int: %d",
				args:     []interface{}{1234567890},
			},
			expected: "info\tsome int: 1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Infof(tt.args.ctx, tt.args.template, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_Warn(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  context.Background(),
				args: []interface{}{"1", "2", "3"},
			},
			expected: "warn\t123\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				args: []interface{}{"some string", 1234567890},
			},
			expected: "warn\tsome string1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Warn(tt.args.ctx, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_Warnf(t *testing.T) {
	type args struct {
		ctx      context.Context
		template string
		args     []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      context.Background(),
				template: "some value: %s",
				args:     []interface{}{"value"},
			},
			expected: "warn\tsome value: value\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				template: "some int: %d",
				args:     []interface{}{1234567890},
			},
			expected: "warn\tsome int: 1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Warnf(tt.args.ctx, tt.args.template, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_Error(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  context.Background(),
				args: []interface{}{"1", "2", "3"},
			},
			expected: "error\t123\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:  mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				args: []interface{}{"some string", 1234567890},
			},
			expected: "error\tsome string1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Error(tt.args.ctx, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_Errorf(t *testing.T) {
	type args struct {
		ctx      context.Context
		template string
		args     []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "WithoutTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      context.Background(),
				template: "some value: %s",
				args:     []interface{}{"value"},
			},
			expected: "error\tsome value: value\n",
		},
		{
			name:   "WithTrace",
			logger: mocks.NewMockLogger(),
			args: args{
				ctx:      mocks.NewContextWithMockSpan(context.Background(), 123, 321),
				template: "some int: %d",
				args:     []interface{}{1234567890},
			},
			expected: "error\tsome int: 1234567890\t{\"trace_id\": \"7b000000000000000000000000000000\", \"span_id\": \"4101000000000000\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.Errorf(tt.args.ctx, tt.args.template, tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_TraceID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "WithSpan#1",
			args: args{
				ctx: mocks.NewContextWithMockSpan(context.Background(), 2144414454365, 1),
			},
			want: "5dd20f49f30100000000000000000000",
		},
		{
			name: "WithSpan#2",
			args: args{
				ctx: mocks.NewContextWithMockSpan(context.Background(), 0, 0),
			},
			want: "",
		},
		{
			name: "WithSpan#3",
			args: args{
				ctx: mocks.NewContextWithMockSpan(context.Background(), 9543901873575874897, 1),
			},
			want: "51d1ac3130c072840000000000000000",
		},
		{
			name: "WithoutSpan#1",
			args: args{
				ctx: context.Background(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(zap.S())

			if got := l.TraceID(tt.args.ctx); got != tt.want {
				t.Errorf("TraceID() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestTraceLogger_With(t *testing.T) {
	type args struct {
		obj  interface{}
		args []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "Pass",
			logger: mocks.NewMockLogger(),
			args: args{
				obj:  struct{ Name string }{Name: "some name"},
				args: []interface{}{"1", "2", "3"},
			},
			expected: "debug\t123\t{\"obj\": {\"Name\":\"some name\"}}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.With("obj", tt.args.obj).Debug(context.Background(), tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}

func TestTraceLogger_WithJSON(t *testing.T) {
	type args struct {
		objData []byte
		args    []interface{}
	}
	tests := []struct {
		name     string
		logger   *mocks.MockLogger
		args     args
		expected string
	}{
		{
			name:   "UnmarshallPass",
			logger: mocks.NewMockLogger(),
			args: args{
				objData: []byte(`{"Name":"some name"}`),
				args:    []interface{}{"1", "2", "3"},
			},
			expected: "debug\t123\t{\"obj\": {\"Name\":\"some name\"}}\n",
		},
		{
			name:   "UnmarshallError",
			logger: mocks.NewMockLogger(),
			args: args{
				objData: []byte(`{{"Name":"some name"}`),
				args:    []interface{}{"1", "2", "3"},
			},
			expected: "debug\t123\t{\"obj\": \"unmarshal failed\", \"failed_json\": \"{{\\\"Name\\\":\\\"some name\\\"}\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.logger.SugaredLogger)

			l.WithJSON("obj", tt.args.objData).Debug(context.Background(), tt.args.args...)

			if err := tt.logger.Sync(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expected, tt.logger.String())
		})
	}
}
