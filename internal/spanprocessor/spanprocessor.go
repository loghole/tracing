package spanprocessor

import (
	"context"
	"sync"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Sampled struct {
	processor tracesdk.SpanProcessor
	sampler   tracesdk.Sampler

	traces map[trace.TraceID]*traceWrapper
	mu     sync.Mutex
}

func NewSampled(
	processor tracesdk.SpanProcessor,
	sampler tracesdk.Sampler,
) *Sampled {
	return &Sampled{
		processor: processor,
		sampler:   sampler,
		traces:    make(map[trace.TraceID]*traceWrapper),
	}
}

func (p *Sampled) OnStart(parent context.Context, span tracesdk.ReadWriteSpan) {
	p.onStart(parent, span)
	p.processor.OnStart(parent, span)
}

func (p *Sampled) OnEnd(span tracesdk.ReadOnlySpan) {
	tr, ok := p.getTrace(span.SpanContext().TraceID())
	if !ok {
		return
	}

	spanID := span.SpanContext().SpanID()

	if !tr.isFinished && !tr.isParent(spanID) {
		return
	}

	hasError := tr.hasError()

	if tr.isParent(spanID) {
		tr.isFinished = true

		for _, span := range tr.extractFinishedSpans() {
			p.send(span.context, span.span, hasError)
		}

		return
	}

	if spanWrapper, ok := tr.spans[spanID]; ok {
		p.send(spanWrapper.context, span, hasError)
	}
}

func (p *Sampled) Shutdown(ctx context.Context) error {
	p.flush()

	return p.processor.Shutdown(ctx)
}

func (p *Sampled) ForceFlush(ctx context.Context) error {
	p.flush()

	return p.processor.ForceFlush(ctx)
}

func (p *Sampled) onStart(parent context.Context, span tracesdk.ReadWriteSpan) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spanCtx := span.SpanContext()

	if val, ok := p.traces[spanCtx.TraceID()]; ok {
		val.storeSpan(parent, span)
	} else {
		p.traces[spanCtx.TraceID()] = &traceWrapper{
			parentSpanID: spanCtx.SpanID(),
			spans: map[trace.SpanID]spanWrapper{
				spanCtx.SpanID(): {span: span, context: parent},
			},
		}
	}
}

func (p *Sampled) getTrace(traceID trace.TraceID) (*traceWrapper, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	tr, ok := p.traces[traceID]

	return tr, ok
}

func (p *Sampled) send(
	parent context.Context,
	span tracesdk.ReadOnlySpan,
	hasError bool,
) {
	if hasError {
		p.processor.OnEnd(span)
	} else if sr := p.sampler.ShouldSample(tracesdk.SamplingParameters{
		ParentContext: parent,
		TraceID:       span.SpanContext().TraceID(),
		Name:          span.Name(),
		Kind:          span.SpanKind(),
		Attributes:    span.Attributes(),
		Links:         linksToTraceLinks(span.Links()),
	}); sr.Decision == tracesdk.RecordAndSample {
		p.processor.OnEnd(span)
	}

	p.removeTraces(span.SpanContext().TraceID())
}

func (p *Sampled) removeTraces(traceID trace.TraceID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	tr, ok := p.traces[traceID]
	if !ok {
		return
	}

	if len(tr.spans) == 0 {
		delete(p.traces, traceID)
	}
}

func (p *Sampled) flush() {
	wg := sync.WaitGroup{}

	p.mu.Lock()

	for _, tr := range p.traces {
		hasError := tr.hasError()

		for _, span := range tr.extractFinishedSpans() {
			wg.Add(1)

			go func(span spanWrapper) {
				defer wg.Done()

				p.send(span.context, span.span, hasError)
			}(span)
		}
	}

	p.mu.Unlock()

	wg.Wait()
}

func linksToTraceLinks(links []tracesdk.Link) []trace.Link {
	resp := make([]trace.Link, 0, len(links))

	for _, l := range links {
		resp = append(resp, trace.Link{
			SpanContext: l.SpanContext,
			Attributes:  l.Attributes,
		})
	}

	return resp
}
