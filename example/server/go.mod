module server

go 1.15

require (
	github.com/loghole/tracing v0.0.0
	go.uber.org/zap v1.20.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

replace github.com/loghole/tracing => ../../
