package tracehttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loghole/tracing/mocks"
)

func TestMiddleware_Middleware(t *testing.T) {
	var (
		tracer, recorder = mocks.NewTracerWithRecorder()
		middleware       = NewMiddleware(tracer)
	)

	middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })).ServeHTTP(requestWR(t))

	require.Len(t, recorder.Ended(), 1)
	assert.NotEmpty(t, recorder.Ended()[0].Attributes())

	middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) })).ServeHTTP(requestWR(t))

	require.Len(t, recorder.Ended(), 2)
	assert.NotEmpty(t, recorder.Ended()[1].Attributes())
}

func requestWR(t *testing.T) (*httptest.ResponseRecorder, *http.Request) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080", http.NoBody)
	require.NoError(t, err, "http.NewRequest")

	return httptest.NewRecorder(), req
}
