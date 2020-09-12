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

	"github.com/gadavy/tracing"
	"github.com/gadavy/tracing/logger"
	"github.com/gadavy/tracing/tracehttp"
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

	log := logger.NewTraceLogger(dev.Sugar())

	serverExample := NewServerExample(log)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		log.Info(ctx, "start server")

		return serverExample.ListenAndServe()
	})

	// Exit from app.
	select {
	case <-exit:
		log.Info(ctx, "stopping application")
	case <-ctx.Done():
		log.Error(ctx, "stopping application with error")
	}

	serverExample.Stop()

	if err := errGroup.Wait(); err != nil {
		log.Error(ctx, err)
	}
}

type ServerExample struct {
	server *http.Server
	tracer *tracing.Tracer
	logger logger.Logger
}

func NewServerExample(logger logger.Logger) *ServerExample {
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

	w.Write([]byte("Error"))
	w.WriteHeader(http.StatusInternalServerError)
}
