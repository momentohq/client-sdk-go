package logger

import (
	"log"
	"strings"
	"time"
)

type loggerLevel string

const (
	INFO  loggerLevel = "INFO"
	DEBUG loggerLevel = "DEBUG"
)

type BuiltinMomentoLogger struct {
	loggerName string
	level      loggerLevel
}

func (l *BuiltinMomentoLogger) Info(message string, args ...string) {
	log.Printf("[%s] %s (%s): %s, %s\n", time.RFC3339, l.level, l.loggerName, message, strings.Join(args, ", "))
}

func (l *BuiltinMomentoLogger) Debug(message string, args ...string) {
	log.Printf("[%s] %s (%s): %s, %s\n", time.RFC3339, l.level, l.loggerName, message, strings.Join(args, ", "))
}

type BuiltinMomentoLoggerFactory struct {
}

func NewBuiltinMomentoLoggerFactory() MomentoLoggerFactory {
	return &BuiltinMomentoLoggerFactory{}
}

func (*BuiltinMomentoLoggerFactory) GetLogger(loggerName string) MomentoLogger {
	log.SetFlags(0)
	return &BuiltinMomentoLogger{loggerName: loggerName, level: INFO}
}
