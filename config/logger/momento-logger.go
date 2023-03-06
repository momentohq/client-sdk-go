package logger

type loggerLevel int

const (
	TRACE loggerLevel = 5
	DEBUG loggerLevel = 10
	INFO  loggerLevel = 20
	WARN  loggerLevel = 30
	ERROR loggerLevel = 40
)

type MomentoLogger interface {
	Trace(message string, args ...string)
	Debug(message string, args ...string)
	Info(message string, args ...string)
	Warn(message string, args ...string)
	Error(message string, args ...string)
}

type MomentoLoggerFactory interface {
	GetLogger(loggerName string, logLevel loggerLevel) MomentoLogger
}
