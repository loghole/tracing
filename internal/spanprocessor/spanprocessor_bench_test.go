package spanprocessor

import (
	"context"
	"math/rand"
	"testing"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func BenchmarkSpanProcessor(b *testing.B) {
	b.Run("Sampled", func(b *testing.B) {
		var (
			processor = NewSampled(NoopSpanProcessor{}, tracesdk.NeverSample())
			tracer    = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor)).Tracer("")
		)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(parallelBenchmark(tracer))
	})

	b.Run("tracesdk.BatchSpanProcessor", func(b *testing.B) {
		var (
			processor = tracesdk.NewBatchSpanProcessor(tracetest.NewNoopExporter())
			tracer    = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor)).Tracer("")
		)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(parallelBenchmark(tracer))
	})

	b.Run("tracesdk.SimpleSpanProcessor", func(b *testing.B) {
		var (
			processor = tracesdk.NewSimpleSpanProcessor(tracetest.NewNoopExporter())
			tracer    = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor)).Tracer("")
		)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(parallelBenchmark(tracer))
	})
}

func parallelBenchmark(tracer trace.Tracer) func(pb *testing.PB) {
	ctx := context.Background()

	funcs := []func(ctx context.Context, tracer trace.Tracer){
		test1,
		test2,
		test3,
	}

	return func(pb *testing.PB) {
		fn := funcs[rand.Intn(len(funcs))]

		for pb.Next() {
			fn(ctx, tracer)
		}
	}
}

func test1(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "span1")
	defer span.End()

	{
		ctx, span := tracer.Start(ctx, "span2")
		defer span.End()

		{
			_, span := tracer.Start(ctx, "span3")
			defer span.End()
		}
	}

	go func() {
		{
			_, span := tracer.Start(ctx, "span4")
			defer span.End()
		}
	}()

	{
		ctx, span := tracer.Start(ctx, "span5")
		defer span.End()

		{
			_, span := tracer.Start(ctx, "span6")
			defer span.End()
		}
	}
}

func test2(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "span1")
	defer span.End()

	{
		_, span := tracer.Start(ctx, "span2")
		defer span.End()
	}
	{
		_, span := tracer.Start(ctx, "span3")
		defer span.End()
	}
	{
		_, span := tracer.Start(ctx, "span4")
		defer span.End()
	}
}

func test3(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "span1")
	defer span.End()

	{
		ctx, span := tracer.Start(ctx, "span2")
		defer span.End()

		{
			ctx, span := tracer.Start(ctx, "span3")
			defer span.End()

			{
				ctx, span := tracer.Start(ctx, "span4")
				defer span.End()

				{
					ctx, span := tracer.Start(ctx, "span5")
					defer span.End()

					{
						_, span := tracer.Start(ctx, "span6")
						defer span.End()
					}
				}
			}
		}
	}
}
