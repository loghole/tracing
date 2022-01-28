package tracehttp

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/loghole/tracing"
	"github.com/loghole/tracing/internal/metrics"
)

const ComponentName = "net/http"

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
	tracer  *tracing.Tracer
	options *Options
}

func NewMiddleware(tracer *tracing.Tracer, options ...Option) *Middleware {
	middleware := &Middleware{tracer: tracer, options: &Options{}}

	for _, option := range options {
		option(middleware.options)
	}

	if middleware.options.Filter == nil {
		middleware.options.Filter = middleware.defaultFilterFunc
	}

	if middleware.options.NameFunc == nil {
		middleware.options.NameFunc = middleware.defaultNameFunc
	}

	return middleware
}

func Handler(tracer *tracing.Tracer, options ...Option) func(next http.Handler) http.Handler {
	m := NewMiddleware(tracer, options...)

	return m.Middleware
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.options.Filter(r) {
			next.ServeHTTP(w, r)

			return
		}

		ctx, span := m.tracer.
			NewSpan().
			WithName(m.options.NameFunc(r)).
			ExtractHeaders(r.Header).
			StartWithContext(r.Context())
		defer span.End()

		span.SetAttributes(
			semconv.HTTPMethodKey.String(r.Method),
			semconv.HTTPURLKey.String(r.URL.String()),
			semconv.HTTPSchemeKey.String(r.URL.Scheme),
			semconv.HTTPRequestContentLengthKey.Int64(r.ContentLength),
		)

		tracker := NewStatusCodeTracker(w)

		next.ServeHTTP(tracker.Writer(), r.WithContext(ctx))

		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(tracker.status))

		if tracker.status >= http.StatusBadRequest {
			metrics.HTTPFailedInputReqCounter.Inc()
			span.SetAttributes(attribute.Bool("error", true))
		} else {
			metrics.HTTPSuccessInputReqCounter.Inc()
		}
	})
}

func (m *Middleware) defaultNameFunc(r *http.Request) string {
	return strings.Join([]string{"HTTP", r.Method, r.RequestURI}, " ")
}

func (m *Middleware) defaultFilterFunc(*http.Request) bool {
	return true
}
