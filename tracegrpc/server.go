package tracegrpc

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	optlog "github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/loghole/tracing/internal/metrics"
	"github.com/loghole/tracing/internal/otgrpc"
)

const componentName = "net/grpc"

// OpenTracingServerInterceptor returns opentracing grpc interceptor.
func UnaryServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		spanContext, _ := extractSpanContext(ctx, tracer)

		span := tracer.StartSpan(defaultNameFunc(info), ext.RPCServerOption(spanContext))
		defer span.Finish()

		ext.Component.Set(span, componentName)

		ctx = opentracing.ContextWithSpan(ctx, span)

		resp, err = handler(ctx, req)
		if err != nil {
			metrics.GRPCFailedInputReqCounter.Inc()
			otgrpc.SetSpanTags(span, err, false)
			span.LogFields(optlog.String("event", "error"), optlog.String("message", err.Error()))
		} else {
			metrics.GRPCSuccessInputReqCounter.Inc()
		}

		return resp, err
	}
}

func StreamServerInterceptor(tracer opentracing.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		spanContext, _ := extractSpanContext(ss.Context(), tracer)

		serverSpan := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(spanContext))
		defer serverSpan.Finish()

		ext.Component.Set(serverSpan, componentName)

		ss = &openTracingServerStream{
			ServerStream: ss,
			ctx:          opentracing.ContextWithSpan(ss.Context(), serverSpan),
		}

		err := handler(srv, ss)
		if err != nil {
			metrics.GRPCFailedInputReqCounter.Inc()
			otgrpc.SetSpanTags(serverSpan, err, false)
			serverSpan.LogFields(optlog.String("event", "error"), optlog.String("message", err.Error()))
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

func extractSpanContext(ctx context.Context, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	return tracer.Extract(opentracing.HTTPHeaders, metadataReaderWriter{md})
}

type metadataReaderWriter struct {
	metadata.MD
}

func (w metadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)

	w.MD[key] = append(w.MD[key], val)
}

func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func defaultNameFunc(r *grpc.UnaryServerInfo) string {
	return strings.Join([]string{"GRPC", r.FullMethod}, " ")
}
