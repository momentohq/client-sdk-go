package logger

import (
	"log"
	"strings"
	"time"
)

type BuiltinMomentoLogger struct {
	loggerName string
	level      LogLevel
}

func (l *BuiltinMomentoLogger) Trace(message string, args ...string) {
	if l.level >= TRACE {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *BuiltinMomentoLogger) Debug(message string, args ...string) {
	if l.level >= DEBUG {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *BuiltinMomentoLogger) Info(message string, args ...string) {
	if l.level >= INFO {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *BuiltinMomentoLogger) Warn(message string, args ...string) {
	if l.level >= WARN {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func (l *BuiltinMomentoLogger) Error(message string, args ...string) {
	if l.level >= ERROR {
		momentoLog(l.level, l.loggerName, message, args...)
	}
}

func momentoLog(level LogLevel, loggerName string, message string, args ...string) {
	if len(args) > 0 {
		log.Printf("[%s] %d (%s): %s, %s\n", time.RFC3339, level, loggerName, message, strings.Join(args, ", "))
	} else {
		log.Printf("[%s] %d (%s): %s\n", time.RFC3339, level, loggerName, message)
	}
}

type BuiltinMomentoLoggerFactory struct {
	level LogLevel
}

func NewBuiltinMomentoLoggerFactory(level LogLevel) MomentoLoggerFactory {
	return &BuiltinMomentoLoggerFactory{
		level: level,
	}
}

func (lf BuiltinMomentoLoggerFactory) GetLogger(loggerName string) MomentoLogger {
	log.SetFlags(0)
	return &BuiltinMomentoLogger{loggerName: loggerName, level: lf.level}
}
