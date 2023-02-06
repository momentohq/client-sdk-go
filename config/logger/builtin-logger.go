package logger

import (
	"log"
)

type BuiltinMomentoLogger struct {
	loggerName string
}

func (l *BuiltinMomentoLogger) Info(message string, args ...any) {
	if args == nil {
		log.Printf(`{"level": "INFO", "message": "%s", "name": "%s"}`, message, l.loggerName)
	} else {
		log.Printf(`{"level": "INFO", "message": "%s", "name": "%s", "%v"}`, message, l.loggerName, args)
	}
}

func (l *BuiltinMomentoLogger) Debug(message string, args ...any) {
	if args == nil {
		log.Printf(`{"level": "DEBUG", "message": "%s", "name": "%s"}`, message, l.loggerName)
	} else {
		log.Printf(`{"level": "DEBUG", "message": "%s", "name": "%s", "%v"}`, message, l.loggerName, args)
	}
}

type BuiltinMomentoLoggerFactory struct {
}

func NewBuiltinMomentoLoggerFactory() MomentoLoggerFactory {
	return &BuiltinMomentoLoggerFactory{}
}

func (*BuiltinMomentoLoggerFactory) GetLogger(loggerName string) MomentoLogger {
	log.SetFlags(0)
	return &BuiltinMomentoLogger{loggerName}
}
