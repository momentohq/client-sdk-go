package logger

type MomentoLogger interface {
	Trace(message string, args ...any)
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
}

type MomentoLoggerFactory interface {
	GetLogger(loggerName string) MomentoLogger
}
