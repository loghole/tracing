package tracehttp

import (
	"io"
	"net/http"
	"testing"
)

func TestStatusCodeTracker_Writer(t *testing.T) {
	// combination 1/32
	{
		inner := struct {
			http.ResponseWriter
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}

	}

	// combination 2/32
	{
		inner := struct {
			http.ResponseWriter
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}

	}

	// combination 3/32
	{
		inner := struct {
			http.ResponseWriter
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}

	}

	// combination 4/32
	{
		inner := struct {
			http.ResponseWriter
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}

	}

	// combination 5/32
	{
		inner := struct {
			http.ResponseWriter
			http.Hijacker
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}

	}

	// combination 6/32
	{
		inner := struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}

	}

	// combination 7/32
	{
		inner := struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}

	}

	// combination 8/32
	{
		inner := struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}

	}

	// combination 9/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}

	}

	// combination 10/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}

	}

	// combination 11/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}

	}

	// combination 12/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 13/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 14/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 15/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 16/32
	{
		inner := struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 17/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 18/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 19/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 20/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 21/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 22/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 23/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 24/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 25/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 26/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 27/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 28/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 29/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 30/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != false {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}

	// combination 31/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != false {
			t.Error("unexpected interface")
		}
	}

	// combination 32/32
	{
		inner := struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
			http.Pusher
		}{}
		w := NewStatusCodeTracker(inner).Writer()
		if _, ok := w.(http.ResponseWriter); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Flusher); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.CloseNotifier); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Hijacker); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(io.ReaderFrom); ok != true {
			t.Error("unexpected interface")
		}
		if _, ok := w.(http.Pusher); ok != true {
			t.Error("unexpected interface")
		}
	}
}
