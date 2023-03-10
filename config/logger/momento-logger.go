package logger

type MomentoLogger interface {
	Trace(message string, args ...string)
	Debug(message string, args ...string)
	Info(message string, args ...string)
	Warn(message string, args ...string)
	Error(message string, args ...string)
}

type MomentoLoggerFactory interface {
	GetLogger(loggerName string) MomentoLogger
}
