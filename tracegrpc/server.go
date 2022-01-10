package tracegrpc

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/loghole/tracing"
	"github.com/loghole/tracing/internal/metrics"
)

const _componentName = "net/grpc"

// UnaryServerInterceptor returns trace grpc interceptor.
func UnaryServerInterceptor(tracer *tracing.Tracer) grpc.UnaryServerInterceptor {
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

		ctx, span := tracer.
			NewSpan().
			WithName(defaultNameFunc(info.FullMethod)).
			ExtractHeaders(http.Header(md)).
			StartWithContext(ctx)
		defer span.End()

		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			attribute.String("component", _componentName),
		)

		resp, err = handler(ctx, req)
		if err != nil {
			metrics.GRPCFailedInputReqCounter.Inc()

			if st, ok := status.FromError(err); ok {
				span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(st.Code())))
			}

			span.SetAttributes(attribute.Bool("error", true))
		} else {
			metrics.GRPCSuccessInputReqCounter.Inc()
		}

		return resp, err
	}
}

func StreamServerInterceptor(tracer *tracing.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}

		ctx, span := tracer.
			NewSpan().
			WithName(defaultNameFunc(info.FullMethod)).
			ExtractHeaders(http.Header(md)).
			StartWithContext(ss.Context())
		defer span.End()

		ss = &openTracingServerStream{ServerStream: ss, ctx: ctx}

		err := handler(srv, ss)
		if err != nil {
			metrics.GRPCFailedInputReqCounter.Inc()

			if st, ok := status.FromError(err); ok {
				span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(st.Code())))
			}

			span.SetAttributes(attribute.Bool("error", true))
		} else {
			metrics.GRPCSuccessInputReqCounter.Inc()
		}

		return err
	}
}

type openTracingServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *openTracingServerStream) Context() context.Context {
	return ss.ctx
}

func defaultNameFunc(method string) string {
	return "GRPC " + method
}
