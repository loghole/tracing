package tracing

import (
	"log"
	"os"
)

type WarnLogger interface {
	Warnf(template string, args ...interface{})
}

//nolint:gochecknoglobals // logger for search bug
var warnLogger WarnLogger

func InitWarnLogger(logger WarnLogger) {
	warnLogger = logger
}

func warnf(template string, args ...interface{}) {
	switch warnLogger != nil {
	case true:
		warnLogger.Warnf(template, args...)
	default:
		log.New(os.Stdout, "tracing: [warning] ", log.Ldate).Printf(template, args...)
	}
}
