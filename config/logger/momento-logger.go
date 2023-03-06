package logger

type loggerLevel string

const (
	INFO  loggerLevel = "INFO"
	DEBUG loggerLevel = "DEBUG"
	WARN  loggerLevel = "WARN"
	ERROR loggerLevel = "ERROR"
	TRACE loggerLevel = "TRACE"
)

type MomentoLogger interface {
	Info(message string, args ...string)
	Debug(message string, args ...string)
	Warn(message string, args ...string)
	Error(message string, args ...string)
	Trace(message string, args ...string)
}

type MomentoLoggerFactory interface {
	GetLogger(loggerName string, logLevel loggerLevel) MomentoLogger
}
