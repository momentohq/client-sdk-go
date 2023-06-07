package momento_default_logger

import (
	"fmt"
	"log"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type LogLevel int

const (
	TRACE LogLevel = 5
	DEBUG LogLevel = 10
	INFO  LogLevel = 20
	WARN  LogLevel = 30
	ERROR LogLevel = 40
)

type DefaultMomentoLogger struct {
	loggerName string
	level      LogLevel
}

func (l *DefaultMomentoLogger) Trace(message string, args ...string) {
	if l.level >= TRACE {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *DefaultMomentoLogger) Debug(message string, args ...string) {
	if l.level >= DEBUG {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *DefaultMomentoLogger) Info(message string, args ...string) {
	if l.level >= INFO {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *DefaultMomentoLogger) Warn(message string, args ...string) {
	if l.level >= WARN {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *DefaultMomentoLogger) Error(message string, args ...string) {
	if l.level >= ERROR {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func momentoLog(level LogLevel, loggerName string, message string, args ...string) {
	anyArgs := make([]any, len(args))
	for i, v := range args {
		anyArgs[i] = v
	}
	finalMessage := fmt.Sprintf(message, anyArgs...)
	log.Printf("[%s] %d (%s): %s\n", time.Now().UTC().Format(time.RFC3339), level, loggerName, finalMessage)
}

type DefaultMomentoLoggerFactory struct {
	level LogLevel
}

func NewDefaultMomentoLoggerFactory(level LogLevel) logger.MomentoLoggerFactory {
	return &DefaultMomentoLoggerFactory{
		level: level,
	}
}

func (lf DefaultMomentoLoggerFactory) GetLogger(loggerName string) logger.MomentoLogger {
	log.SetFlags(0)
	return &DefaultMomentoLogger{loggerName: loggerName, level: lf.level}
}
