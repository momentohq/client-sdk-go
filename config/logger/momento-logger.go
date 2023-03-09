package logger

type LogLevel int

const (
	TRACE LogLevel = 5
	DEBUG LogLevel = 10
	INFO  LogLevel = 20
	WARN  LogLevel = 30
	ERROR LogLevel = 40
)

type MomentoLogger interface {
	Trace(message string, args ...string)
	Debug(message string, args ...string)
	Info(message string, args ...string)
	Warn(message string, args ...string)
	Error(message string, args ...string)
}

type MomentoLoggerFactory interface {
	GetLogger(loggerName string, logLevel LogLevel) MomentoLogger
}
