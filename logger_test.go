package tracing

import (
	"context"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestTraceLogger_argsWithAction(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		args     []interface{}
		expected []interface{}
	}{
		{
			name:     "WithAction",
			ctx:      ContextWithAction(context.Background(), "some action"),
			args:     []interface{}{1, "aaa", "bbb", "ccc"},
			expected: []interface{}{"action=some action; ", 1, "aaa", "bbb", "ccc"},
		},
		{
			name:     "WithoutAction",
			ctx:      context.Background(),
			args:     []interface{}{1, "aaa", "bbb", "ccc"},
			expected: []interface{}{1, "aaa", "bbb", "ccc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := argsWithAction(tt.ctx, tt.args)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTraceLogger_templateWithAction(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		template string
		expected string
	}{
		{
			name:     "WithAction",
			ctx:      ContextWithAction(context.Background(), "some action"),
			template: "some template: %v, %s",
			expected: "action=some action; some template: %v, %s",
		},
		{
			name:     "WithoutAction",
			ctx:      context.Background(),
			template: "some template: %v, %s",
			expected: "some template: %v, %s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := templateWithAction(tt.ctx, tt.template)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTraceLogger_withAction(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "WithAction",
			ctx:      ContextWithAction(context.Background(), "some action"),
			expected: "action=some action; ",
		},
		{
			name:     "WithoutAction",
			ctx:      context.Background(),
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := withAction(tt.ctx)

			assert.Equal(t, tt.expected, result)
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

func TestTraceLogger_Debug(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		input    []interface{}
	}{
		{
			name:     "#1",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"action=1; ", "aaa", 222, "qwerty"}},
			ctx:      ContextWithAction(context.Background(), "1"),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"aaa", 222, "qwerty"}},
			ctx:      context.Background(),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Debug", tt.loggerIn...)

			logger := TraceLogger{logger: tt.logger}

			logger.Debug(tt.ctx, tt.input...)
		})
	}
}

func TestTraceLogger_Debugf(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		template string
		args     []interface{}
	}{
		{
			name:     "#1",
			ctx:      ContextWithAction(context.Background(), "1"),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"action=1; some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			ctx:      context.Background(),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Debugf", tt.loggerIn...)

			logger := NewTraceLogger(tt.logger)

			logger.Debugf(tt.ctx, tt.template, tt.args...)
		})
	}
}

func TestTraceLogger_Info(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		input    []interface{}
	}{
		{
			name:     "#1",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"action=1; ", "aaa", 222, "qwerty"}},
			ctx:      ContextWithAction(context.Background(), "1"),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"aaa", 222, "qwerty"}},
			ctx:      context.Background(),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Info", tt.loggerIn...)

			logger := NewTraceLogger(tt.logger)

			logger.Info(tt.ctx, tt.input...)
		})
	}
}

func TestTraceLogger_Infof(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		template string
		args     []interface{}
	}{
		{
			name:     "#1",
			ctx:      ContextWithAction(context.Background(), "1"),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"action=1; some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			ctx:      context.Background(),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Infof", tt.loggerIn...)

			logger := NewTraceLogger(tt.logger)

			logger.Infof(tt.ctx, tt.template, tt.args...)
		})
	}
}

func TestTraceLogger_Warn(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		input    []interface{}
	}{
		{
			name:     "#1",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"action=1; ", "aaa", 222, "qwerty"}},
			ctx:      ContextWithAction(context.Background(), "1"),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"aaa", 222, "qwerty"}},
			ctx:      context.Background(),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Warn", tt.loggerIn...)

			logger := NewTraceLogger(tt.logger)

			logger.Warn(tt.ctx, tt.input...)
		})
	}
}

func TestTraceLogger_Warnf(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		template string
		args     []interface{}
	}{
		{
			name:     "#1",
			ctx:      ContextWithAction(context.Background(), "1"),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"action=1; some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			ctx:      context.Background(),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Warnf", tt.loggerIn...)

			logger := NewTraceLogger(tt.logger)

			logger.Warnf(tt.ctx, tt.template, tt.args...)
		})
	}
}

func TestTraceLogger_Error(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		input    []interface{}
	}{
		{
			name:     "#1",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"action=1; ", "aaa", 222, "qwerty"}},
			ctx:      opentracing.ContextWithSpan(ContextWithAction(context.Background(), "1"), mocktracer.New().StartSpan("test")),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{[]interface{}{"aaa", 222, "qwerty"}},
			ctx:      context.Background(),
			input:    []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Error", tt.loggerIn...)

			logger := NewTraceLogger(tt.logger)

			logger.Error(tt.ctx, tt.input...)

			if span := opentracing.SpanFromContext(tt.ctx); span != nil {
				mockSpan, ok := span.(*mocktracer.MockSpan)
				if !ok {
					t.Fatal("expected mocktracer.MockSpan")
				}

				assert.Equal(t, true, mockSpan.Tag("error"))
			}
		})
	}
}

func TestTraceLogger_Errorf(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		logger   *MockBaseLogger
		loggerIn []interface{}
		template string
		args     []interface{}
	}{
		{
			name:     "#1",
			ctx:      opentracing.ContextWithSpan(ContextWithAction(context.Background(), "1"), mocktracer.New().StartSpan("test")),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"action=1; some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
		{
			name:     "#2",
			ctx:      context.Background(),
			logger:   &MockBaseLogger{},
			loggerIn: []interface{}{"some template: %s %s %s %s", []interface{}{"aaa", 222, "qwerty"}},
			template: "some template: %s %s %s %s",
			args:     []interface{}{"aaa", 222, "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.On("Errorf", tt.loggerIn...)

			logger := NewTraceLogger(tt.logger)

			logger.Errorf(tt.ctx, tt.template, tt.args...)

			if span := opentracing.SpanFromContext(tt.ctx); span != nil {
				mockSpan, ok := span.(*mocktracer.MockSpan)
				if !ok {
					t.Fatal("expected mocktracer.MockSpan")
				}

				assert.Equal(t, true, mockSpan.Tag("error"))
			}
		})
	}
}

type MockBaseLogger struct {
	mock.Mock
}

func (m *MockBaseLogger) Debug(args ...interface{}) {
	m.Called(args)
}
func (m *MockBaseLogger) Debugf(template string, args ...interface{}) {
	m.Called(template, args)
}
func (m *MockBaseLogger) Info(args ...interface{}) {
	m.Called(args)
}
func (m *MockBaseLogger) Infof(template string, args ...interface{}) {
	m.Called(template, args)
}
func (m *MockBaseLogger) Warn(args ...interface{}) {
	m.Called(args)
}
func (m *MockBaseLogger) Warnf(template string, args ...interface{}) {
	m.Called(template, args)
}
func (m *MockBaseLogger) Error(args ...interface{}) {
	m.Called(args)
}
func (m *MockBaseLogger) Errorf(template string, args ...interface{}) {
	m.Called(template, args)
}
