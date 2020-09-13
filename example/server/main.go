package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracehttp"
	"github.com/loghole/tracing/tracelog"
)

const (
	jaegerURL  = "127.0.0.1:6831"
	serverAddr = "127.0.0.1:38572"
)

func main() {
	dev, err := zap.NewProduction(zap.AddStacktrace(zap.FatalLevel))
	if err != nil {
		panic(err)
	}

	logger := tracelog.NewTraceLogger(dev.Sugar())

	serverExample := NewServerExample(logger)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		logger.Info(ctx, "start server")

		return serverExample.ListenAndServe()
	})

	// Exit from app.
	select {
	case <-exit:
		logger.Info(ctx, "stopping application")
	case <-ctx.Done():
		logger.Error(ctx, "stopping application with error")
	}

	serverExample.Stop()

	if err := errGroup.Wait(); err != nil {
		logger.Error(ctx, err)
	}
}

type ServerExample struct {
	server *http.Server
	tracer *tracing.Tracer
	logger tracelog.Logger
}

func NewServerExample(logger tracelog.Logger) *ServerExample {
	tracer, err := tracing.NewTracer(tracing.DefaultConfiguration("example_server", jaegerURL))
	if err != nil {
		panic(err)
	}

	return &ServerExample{
		tracer: tracer,
		logger: logger,
	}
}

func (e *ServerExample) ListenAndServe() error {
	routes := http.NewServeMux()
	routes.HandleFunc("/success", e.HandlerSuccess)
	routes.HandleFunc("/error", e.HandlerError)

	middleware := tracehttp.NewMiddleware(e.tracer)

	e.server = &http.Server{Addr: serverAddr, Handler: middleware.Middleware(routes)}

	return e.server.ListenAndServe()
}

func (e *ServerExample) Stop() {
	if err := e.tracer.Close(); err != nil {
		e.logger.Error(context.TODO(), "close tracer: ", err)
	}

	if err := e.server.Shutdown(context.Background()); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			e.logger.Error(context.TODO(), "close server: ", err)
		}
	}
}

func (e *ServerExample) HandlerSuccess(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	defer tracing.ChildSpan(&ctx).Finish()

	time.Sleep(time.Millisecond * 500)

	w.Write([]byte("OK"))
}

func (e *ServerExample) HandlerError(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	defer tracing.ChildSpan(&ctx).Finish()

	time.Sleep(time.Millisecond * 500)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Error"))
}
