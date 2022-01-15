# Tracing

[![GoDoc](https://pkg.go.dev/badge/github.com/loghole/tracing)](https://pkg.go.dev/github.com/loghole/tracing)
[![Go Report Card](https://goreportcard.com/badge/github.com/loghole/tracing)](https://goreportcard.com/report/github.com/loghole/tracing)

Tracing is a [go.opentelemetry.io](https://github.com/open-telemetry/opentelemetry-go) Tracer wrapper for instrumenting
golang applications to collect traces in jaeger.

# Install

```sh
go get github.com/loghole/tracing
```

# Usage

```go
package main

import (
	"context"
	"time"

	"github.com/loghole/tracing"
)

func main() {
	tracer, err := tracing.NewTracer(tracing.DefaultConfiguration("example", "udp://127.0.0.1:6831"))
	if err != nil {
		panic(err)
	}

	defer tracer.Close()

	ctx, span := tracer.NewSpan().WithName("root").StartWithContext(context.Background())
	defer span.End()

	next(ctx)
}

func next(ctx context.Context) {
	defer tracing.ChildSpan(&ctx).End()

	// Some work...
	time.Sleep(time.Second)
}
```

# Examples

- [HTTP Client](https://github.com/loghole/tracing/blob/85a206d9aa6242f693283e159ac428dc23ea9c99/example/client/main.go)
- [HTTP Server](https://github.com/loghole/tracing/blob/85a206d9aa6242f693283e159ac428dc23ea9c99/example/server/main.go)
- [Simple](https://github.com/loghole/tracing/blob/85a206d9aa6242f693283e159ac428dc23ea9c99/example/simple/main.go)
