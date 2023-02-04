package logger

type MomentoLogger interface {
	Info(message string, args ...any)
	Debug(message string, args ...any)
}

type MomentoLoggerFactory interface {
	GetLogger(loggerName string) MomentoLogger
}
