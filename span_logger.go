package tracing

type WarnLogger interface {
	Warnf(template string, args ...interface{})
}

//nolint:gochecknoglobals // logger for search bug
var warnLogger WarnLogger

func InitWarnLogger(logger WarnLogger) {
	warnLogger = logger
}
