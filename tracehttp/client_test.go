package tracehttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loghole/tracing/mocks"
)

func TestClient_Do(t *testing.T) {
	ctx := context.Background()

	tracer, recorder := mocks.NewTracerWithRecorder()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	}))

	client := NewClient(tracer, server.Client())

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, http.NoBody)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmptyf(t, req.Header.Get("Traceparent"), "empty tracer header")

	// error
	req, err = http.NewRequestWithContext(ctx, "GET", "/", http.NoBody)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.Error(t, err)

	ended := recorder.Ended()
	assert.Len(t, ended, 2)
}

func TestClient_Get(t *testing.T) {
	ctx := context.Background()

	tracer, recorder := mocks.NewTracerWithRecorder()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	}))

	client := NewClient(tracer, server.Client())

	resp, err := client.Get(ctx, server.URL)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// error
	_, err = client.Get(ctx, "/")
	require.Error(t, err)

	ended := recorder.Ended()
	assert.Len(t, ended, 2)
}

func TestClient_Post(t *testing.T) {
	ctx := context.Background()

	tracer, recorder := mocks.NewTracerWithRecorder()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	}))

	client := NewClient(tracer, server.Client())

	resp, err := client.Post(ctx, server.URL, "application/json", bytes.NewReader([]byte("{}")))
	require.NoError(t, err)
	require.NotNil(t, resp)

	// error
	_, err = client.Post(ctx, "/", "application/json", bytes.NewReader([]byte("{}")))
	require.Error(t, err)

	ended := recorder.Ended()
	assert.Len(t, ended, 2)
}

func TestClient_PostForm(t *testing.T) {
	ctx := context.Background()

	tracer, recorder := mocks.NewTracerWithRecorder()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	}))

	client := NewClient(tracer, server.Client())

	resp, err := client.PostForm(ctx, server.URL, url.Values{"key": []string{"value"}})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// error
	_, err = client.PostForm(ctx, "/", url.Values{"key": []string{"value"}})
	require.Error(t, err)

	ended := recorder.Ended()
	assert.Len(t, ended, 2)
}
