package logger

import (
	"fmt"
	"log"
)

type BuitlinMomentoLogger struct {
	loggerName string
}

func (l *BuitlinMomentoLogger) Info(message string, args ...any) {
	str := ""
	if args == nil {
		str = fmt.Sprintf(`{"level": "INFO", "message": "%s", "name": "%s"}`, message, l.loggerName)
	} else {
		str = fmt.Sprintf(`{"level": "INFO", "message": "%s", "name": "%s", "%v"}`, message, l.loggerName, args)
	}
	log.SetFlags(0)
	log.Print(str)
}

func (l *BuitlinMomentoLogger) Debug(message string, args ...any) {
	str := ""
	if args == nil {
		str = fmt.Sprintf(`{"level": "DEBUG", "message": "%s", "name": "%s"}`, message, l.loggerName)
	} else {
		str = fmt.Sprintf(`{"level": "DEBUG", "message": "%s", "name": "%s", "%v"}`, message, l.loggerName, args)
	}
	log.SetFlags(0)
	log.Print(str)
}

type BuiltinMomentoLoggerFactory struct {
}

func NewBuitlinMomentoLoggerFactory() MomentoLoggerFactory {
	return &BuiltinMomentoLoggerFactory{}
}

func (*BuiltinMomentoLoggerFactory) GetLogger(loggerName string) MomentoLogger {
	return &BuitlinMomentoLogger{loggerName}
}
