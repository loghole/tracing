package tracing

import (
	"context"
	"time"
)

type noCancelCtx struct {
	parent context.Context // nolint:containedctx // need parent context.
}

func (c noCancelCtx) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c noCancelCtx) Done() <-chan struct{}             { return nil }
func (c noCancelCtx) Err() error                        { return nil }
func (c noCancelCtx) Value(key interface{}) interface{} { return c.parent.Value(key) }

// ContextWithoutCancel returns a context that is never canceled.
func ContextWithoutCancel(parent context.Context) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	return noCancelCtx{parent: parent}
}
