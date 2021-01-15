package tracegrpc

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/loghole/tracing/internal/metrics"
	"github.com/loghole/tracing/internal/otgrpc"
)

const loadBalancing = `{"loadBalancingPolicy":"round_robin","loadBalancingConfig":[{"round_robin":{}}]}`

func Dial(target string, tracer opentracing.Tracer, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), target, tracer, opts...)
}

func DialContext(
	ctx context.Context,
	target string,
	tracer opentracing.Tracer,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	// Init default options.
	opts = append(opts,
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(loadBalancing),
		grpc.WithChainUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer), metricUnaryInterceptor()),
		grpc.WithChainStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer), metricStreamInterceptor()),
	)

	return grpc.DialContext(ctx, target, opts...)
}

func metricUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			metrics.GRPCFailedOutputReqCounter.Inc()
		} else {
			metrics.GRPCSuccessOutputReqCounter.Inc()
		}

		return err
	}
}

func metricStreamInterceptor() grpc.StreamClientInterceptor {
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
