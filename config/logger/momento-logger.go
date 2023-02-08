package logger

type MomentoLogger interface {
	Info(message string, args ...string)
	Debug(message string, args ...string)
}

type MomentoLoggerFactory interface {
	GetLogger(loggerName string) MomentoLogger
}
