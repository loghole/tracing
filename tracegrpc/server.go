package tracegrpc

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/loghole/tracing/internal/metrics"
)

// UnaryServerInterceptor returns trace grpc interceptor.
func UnaryServerInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ctx = new(propagation.TraceContext).Extract(ctx, propagation.HeaderCarrier(md))

		ctx, span := tracer.Start(ctx, defaultNameFunc(info.FullMethod), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		resp, err = handler(ctx, req)
		if err != nil {
			metrics.GRPCFailedInputReqCounter.Inc()
		} else {
			metrics.GRPCSuccessInputReqCounter.Inc()
		}

		setAttributes(span, info.FullMethod, err)

		return resp, err
	}
}

func StreamServerInterceptor(tracer trace.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}

		ctx := new(propagation.TraceContext).Extract(ss.Context(), propagation.HeaderCarrier(md))

		ctx, span := tracer.Start(ctx, defaultNameFunc(info.FullMethod), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		ss = &tracingServerStream{ServerStream: ss, ctx: ctx}

		err := handler(srv, ss)
		if err != nil {
			metrics.GRPCFailedInputReqCounter.Inc()
		} else {
			metrics.GRPCSuccessInputReqCounter.Inc()
		}

		setAttributes(span, info.FullMethod, err)

		return err
	}
}

type tracingServerStream struct {
	grpc.ServerStream
	ctx context.Context // nolint:containedctx // need internal context.
}

func (ss *tracingServerStream) Context() context.Context {
	return ss.ctx
}

func defaultNameFunc(method string) string {
	return "GRPC " + method
}
