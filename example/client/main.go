package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
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

	client := NewClientExample(logger)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		logger.Info(ctx, "start server")

		return client.Run()
	})

	// Exit from app.
	select {
	case <-exit:
		logger.Info(ctx, "stopping application")
	case <-ctx.Done():
		logger.Error(ctx, "stopping application with error")
	}

	client.Stop()

	if err := errGroup.Wait(); err != nil {
		logger.Error(ctx, err)
	}
}

type ClientExample struct {
	client *tracehttp.Client
	tracer *tracing.Tracer
	logger tracelog.Logger

	close chan struct{}
	rnd   *rand.Rand
}

func NewClientExample(logger tracelog.Logger) *ClientExample {
	tracer, err := tracing.NewTracer(tracing.DefaultConfiguration("example_client", jaegerURL))
	if err != nil {
		panic(err)
	}

	return &ClientExample{
		client: tracehttp.NewClient(tracer, http.DefaultClient, false),
		tracer: tracer,
		logger: logger,
		close:  make(chan struct{}),
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *ClientExample) Run() error {
	var (
		ticSuccess = time.NewTicker(time.Millisecond * 3500)
		ticError   = time.NewTicker(time.Millisecond * 7500)
		ticGoogle  = time.NewTicker(time.Millisecond * 5500)
	)

	for {
		select {
		case <-e.close:
			return nil
		case <-ticSuccess.C:
			e.SendRequest(fmt.Sprintf("http://%s%s", serverAddr, "/success"))
		case <-ticError.C:
			e.SendRequest(fmt.Sprintf("http://%s%s", serverAddr, "/error"))
		case <-ticGoogle.C:
			e.SendRequest("https://google.com/robots.txt")
		}
	}
}

func (e *ClientExample) Stop() {
	if err := e.tracer.Close(); err != nil {
		e.logger.Error(context.TODO(), "close tracer: ", err)
	}

	close(e.close)
}

func (e *ClientExample) SendRequest(url string) {
	span, ctx := e.tracer.NewSpan().BuildWithContext(context.Background())
	defer span.Finish()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		e.logger.Errorf(ctx, "NewRequestWithContext: %v", err)

		return
	}

	resp, err := e.client.Do(req)
	if err != nil {
		e.logger.Errorf(ctx, "Do request: %v", err)

		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e.logger.Errorf(ctx, "ReadAll: %v", err)

		return
	}

	e.logger.With("text", string(data)).Infof(ctx, "success")
}
