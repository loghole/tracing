package spanprocessor

import (
	"context"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/codes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type wrapper struct {
	parent    tracesdk.ReadWriteSpan
	parentCtx context.Context // nolint:containedctx // need context.

	spans map[trace.SpanID]tracesdk.ReadWriteSpan

	sampled bool
	once    sync.Once

	sync.Mutex
}

func (w *wrapper) isSampled(sampler tracesdk.Sampler) bool {
	w.once.Do(func() {
		w.sampled = w.checkSampled(sampler)
	})

	return w.sampled
}

func (w *wrapper) checkSampled(sampler tracesdk.Sampler) bool {
	for _, span := range w.spans {
		if span.Status().Code == codes.Error {
			return true
		}
	}

	for _, span := range w.spans {
		for _, attr := range span.Attributes() {
			if strings.EqualFold(string(attr.Key), "error") {
				return true
			}
		}
	}

	result := sampler.ShouldSample(tracesdk.SamplingParameters{
		ParentContext: w.parentCtx,
		TraceID:       w.parent.SpanContext().TraceID(),
		Name:          w.parent.Name(),
		Kind:          w.parent.SpanKind(),
		Attributes:    w.parent.Attributes(),
		Links:         nil, // skip links, because they are not used in samplers.
	})

	return result.Decision == tracesdk.RecordAndSample
}

type Sampled struct {
	processor tracesdk.SpanProcessor
	sampler   tracesdk.Sampler

	traces map[trace.TraceID]*wrapper
	mu     sync.RWMutex
}

func NewSampled(
	processor tracesdk.SpanProcessor,
	sampler tracesdk.Sampler,
) *Sampled {
	return &Sampled{
		processor: processor,
		sampler:   sampler,
		traces:    make(map[trace.TraceID]*wrapper),
	}
}

func (p *Sampled) OnStart(parent context.Context, span tracesdk.ReadWriteSpan) {
	if !span.IsRecording() {
		return
	}

	var (
		spanCtx = span.SpanContext()
		traceID = spanCtx.TraceID()
		spanID  = spanCtx.SpanID()
	)

	p.mu.Lock()
	defer p.mu.Unlock()

	wr, ok := p.traces[traceID]
	if ok {
		wr.spans[spanID] = span

		return
	}

	p.traces[traceID] = &wrapper{
		parent:    span,
		parentCtx: parent,
		spans:     map[trace.SpanID]tracesdk.ReadWriteSpan{spanID: span},
	}
}

func (p *Sampled) OnEnd(span tracesdk.ReadOnlySpan) {
	var (
		spanCtx = span.SpanContext()
		traceID = spanCtx.TraceID()
	)

	p.mu.Lock()
	defer p.mu.Unlock()

	wr, ok := p.traces[traceID]
	if !ok {
		return
	}

	if !wr.parent.SpanContext().Equal(spanCtx) && wr.parent.IsRecording() {
		return
	}

	p.finishWrapper(wr)
}

func (p *Sampled) Shutdown(ctx context.Context) error {
	p.flush() // nolint:contextcheck // not need.

	return p.processor.Shutdown(ctx)
}

func (p *Sampled) ForceFlush(ctx context.Context) error {
	p.flush() // nolint:contextcheck // not need.

	return p.processor.ForceFlush(ctx)
}

func (p *Sampled) flush() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, wr := range p.traces {
		p.finishWrapper(wr) // nolint:contextcheck // not need.
	}
}

func (p *Sampled) finishWrapper(wr *wrapper) {
	traceID := wr.parent.SpanContext().TraceID()

	if !wr.isSampled(p.sampler) {
		delete(p.traces, traceID)

		return
	}

	for _, span := range wr.spans {
		if span.EndTime().IsZero() {
			continue
		}

		delete(wr.spans, span.SpanContext().SpanID())

		p.processor.OnStart(trace.ContextWithSpan(context.Background(), span), span)
		p.processor.OnEnd(span)
	}

	if len(wr.spans) == 0 {
		delete(p.traces, traceID)
	}
}
