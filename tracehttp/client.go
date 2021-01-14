package tracehttp

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"

	"github.com/loghole/tracing/internal/metrics"
)

type Client struct {
	client   *http.Client
	tracer   opentracing.Tracer
	extended bool
}

func NewClient(tracer opentracing.Tracer, client *http.Client, extended bool) *Client {
	client.Transport = NewTransport(tracer, client.Transport, extended)

	return &Client{client: client, tracer: tracer, extended: extended}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.extended {
		var ht *nethttp.Tracer

		req, ht = nethttp.TraceRequest(c.tracer, req)

		defer ht.Finish()
	}

	resp, err := c.client.Do(req)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		metrics.HTTPFailedOutputReqCounter.Inc()
	} else {
		metrics.HTTPSuccessOutputReqCounter.Inc()
	}

	return resp, err
}

func (c *Client) Get(ctx context.Context, uri string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Post(ctx context.Context, uri, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

func (c *Client) PostForm(ctx context.Context, uri string, data url.Values) (*http.Response, error) {
	return c.Post(ctx, uri, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
