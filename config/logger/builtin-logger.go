package logger

import (
	"log"
	"strings"
	"time"
)

type BuiltinMomentoLogger struct {
	loggerName string
	level      loggerLevel
}

func (l *BuiltinMomentoLogger) Info(message string, args ...string) {
	if l.level == INFO {
		if len(args) != 0 {
			logWithArgs(l.level, l.loggerName, message, args...)
		} else {
			logWithoutArgs(l.level, l.loggerName, message)
		}
	}
}

func (l *BuiltinMomentoLogger) Debug(message string, args ...string) {
	if l.level == DEBUG {
		if len(args) != 0 {
			logWithArgs(l.level, l.loggerName, message, args...)
		} else {
			logWithoutArgs(l.level, l.loggerName, message)
		}
	}
}

func (l *BuiltinMomentoLogger) Warn(message string, args ...string) {
	if l.level == WARN {
		if len(args) != 0 {
			logWithArgs(l.level, l.loggerName, message, args...)
		} else {
			logWithoutArgs(l.level, l.loggerName, message)
		}
	}
}

func (l *BuiltinMomentoLogger) Error(message string, args ...string) {
	if l.level == ERROR {
		if len(args) != 0 {
			logWithArgs(l.level, l.loggerName, message, args...)
		} else {
			logWithoutArgs(l.level, l.loggerName, message)
		}
	}
}

func (l *BuiltinMomentoLogger) Trace(message string, args ...string) {
	if l.level == TRACE {
		if len(args) != 0 {
			logWithArgs(l.level, l.loggerName, message, args...)
		} else {
			logWithoutArgs(l.level, l.loggerName, message)
		}
	}
}

func logWithArgs(level loggerLevel, loggerName string, message string, args ...string) {
	log.Printf("[%s] %s (%s): %s, %s\n", time.RFC3339, level, loggerName, message, strings.Join(args, ", "))
}

func logWithoutArgs(level loggerLevel, loggerName string, message string) {
	log.Printf("[%s] %s (%s): %s\n", time.RFC3339, level, loggerName, message)
}

type BuiltinMomentoLoggerFactory struct {
}

func NewBuiltinMomentoLoggerFactory() MomentoLoggerFactory {
	return &BuiltinMomentoLoggerFactory{}
}

func (*BuiltinMomentoLoggerFactory) GetLogger(loggerName string, level loggerLevel) MomentoLogger {
	log.SetFlags(0)
	return &BuiltinMomentoLogger{loggerName: loggerName, level: level}
}
