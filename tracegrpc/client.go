package tracegrpc

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/loghole/tracing"
	"github.com/loghole/tracing/internal/metrics"
)

const loadBalancing = `{"loadBalancingPolicy":"round_robin","loadBalancingConfig":[{"round_robin":{}}]}`

func Dial(target string, tracer trace.Tracer, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), target, tracer, opts...)
}

func DialContext(
	ctx context.Context,
	target string,
	tracer trace.Tracer,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	// Init default options.
	opts = append(opts,
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(loadBalancing),
		grpc.WithChainUnaryInterceptor(UnaryClientInterceptor(tracer)),
		grpc.WithChainStreamInterceptor(StreamClientInterceptor()),
	)

	return grpc.DialContext(ctx, target, opts...)
}

func UnaryClientInterceptor(tracer trace.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx, span := tracer.Start(ctx, defaultNameFunc(method), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		tracing.InjectHeaders(ctx, http.Header(md))

		err := invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
		if err != nil {
			metrics.GRPCFailedOutputReqCounter.Inc()
		} else {
			metrics.GRPCSuccessOutputReqCounter.Inc()
		}

		setAttributes(span, method, err)

		return err
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			metrics.GRPCFailedOutputReqCounter.Inc()
		} else {
			metrics.GRPCSuccessOutputReqCounter.Inc()
		}

		return stream, err
	}
}

func setAttributes(span trace.Span, method string, err error) {
	st, _ := status.FromError(err)

	span.SetAttributes(
		semconv.RPCSystemKey.String("GRPC"),
		semconv.RPCMethodKey.String(method),
		semconv.RPCGRPCStatusCodeKey.Int(int(st.Code())),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error")
	}
}
