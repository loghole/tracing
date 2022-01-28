package tracehttp

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/loghole/tracing/internal/metrics"
)

type Option func(options *Options)

func WithNameFunc(f func(r *http.Request) string) Option {
	return func(options *Options) {
		options.NameFunc = f
	}
}

func WithFilterFunc(f func(r *http.Request) bool) Option {
	return func(options *Options) {
		options.Filter = f
	}
}

type Options struct {
	NameFunc func(r *http.Request) string
	Filter   func(r *http.Request) bool
}

type Middleware struct {
	tracer  trace.Tracer
	options *Options
}

func NewMiddleware(tracer trace.Tracer, options ...Option) *Middleware {
	middleware := &Middleware{tracer: tracer, options: &Options{}}

	for _, option := range options {
		option(middleware.options)
	}

	if middleware.options.Filter == nil {
		middleware.options.Filter = middleware.defaultFilterFunc
	}

	if middleware.options.NameFunc == nil {
		middleware.options.NameFunc = defaultNameFunc
	}

	return middleware
}

func Handler(tracer trace.Tracer, options ...Option) func(next http.Handler) http.Handler {
	m := NewMiddleware(tracer, options...)

	return m.Middleware
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.options.Filter(r) {
			next.ServeHTTP(w, r)

			return
		}

		ctx := new(propagation.TraceContext).Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		ctx, span := m.tracer.Start(ctx, defaultNameFunc(r), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(
			semconv.HTTPMethodKey.String(r.Method),
			semconv.HTTPURLKey.String(r.URL.String()),
			semconv.HTTPSchemeKey.String(r.URL.Scheme),
			semconv.HTTPRequestContentLengthKey.Int64(r.ContentLength),
		)

		tracker := NewStatusCodeTracker(w)

		next.ServeHTTP(tracker, r.WithContext(ctx))

		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(tracker.status))

		if tracker.status >= http.StatusBadRequest {
			metrics.HTTPFailedInputReqCounter.Inc()
			span.SetAttributes(attribute.Bool("error", true))
		} else {
			metrics.HTTPSuccessInputReqCounter.Inc()
		}
	})
}

func (m *Middleware) defaultFilterFunc(*http.Request) bool {
	return true
}
