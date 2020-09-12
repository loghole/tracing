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

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"go.uber.org/zap"

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

	client := NewClientExample(log)
	client.Test()
	defer client.Stop()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// errGroup, ctx := errgroup.WithContext(context.Background())
	//
	// errGroup.Go(func() error {
	// 	log.Info(ctx, "start server")
	//
	// 	return client.Run()
	// })
	//
	// // Exit from app.
	// select {
	// case <-exit:
	// 	log.Info(ctx, "stopping application")
	// case <-ctx.Done():
	// 	log.Error(ctx, "stopping application with error")
	// }
	//
	// client.Stop()
	//
	// if err := errGroup.Wait(); err != nil {
	// 	log.Error(ctx, err)
	// }
}

type ClientExample struct {
	client *http.Client
	tracer *tracing.Tracer
	logger logger.Logger

	close chan struct{}
	rnd   *rand.Rand
}

func NewClientExample(logger logger.Logger) *ClientExample {
	tracer, err := tracing.NewTracer(tracing.DefaultConfiguration("example_client", jaegerURL))
	if err != nil {
		panic(err)
	}

	return &ClientExample{
		client: tracehttp.NewClient(tracer, http.DefaultClient),
		tracer: tracer,
		logger: logger,
		close:  make(chan struct{}),
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *ClientExample) Run() error {
	var (
		ticSuccess = time.NewTicker(time.Millisecond * 350)
		ticError   = time.NewTicker(time.Millisecond * 750)
		ticGoogle  = time.NewTicker(time.Millisecond * 2500)
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
			e.SendRequest("https://www.google.com/robots.txt")
		}

		return nil
	}
}

func (e *ClientExample) Stop() {
	if err := e.tracer.Close(); err != nil {
		e.logger.Error(context.TODO(), "close tracer: ", err)
	}

	close(e.close)
}
func (e *ClientExample) Test() error {
	span, ctx := e.tracer.NewSpan().BuildWithContext(context.Background())
	defer span.Finish()

	client := &http.Client{Transport: &nethttp.Transport{}}

	req, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx) // extend existing trace, if any

	req, ht := nethttp.TraceRequest(e.tracer, req)
	defer ht.Finish()

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()

	return nil
}
func (e *ClientExample) SendRequest(url string) {
	span, ctx := e.tracer.NewSpan().BuildWithContext(context.Background())
	defer span.Finish()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		e.logger.Errorf(ctx, "NewRequestWithContext: %v", err)

		return
	}

	resp, err := tracehttp.NewClient(e.tracer, http.DefaultClient).Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		e.logger.Errorf(ctx, "Do request: %v", err)

		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e.logger.Errorf(ctx, "ReadAll: %v", err)

		return
	}

	e.logger.With("text", string(data)).Infof(ctx, "success")
}
