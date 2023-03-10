package momento_logrus

import (
	"log"

	"github.com/momentohq/client-sdk-go/config/logger"
	logrus "github.com/sirupsen/logrus"
)

type LogrusMomentoLogger struct {
	logrusLogger *logrus.Entry
}

func (l LogrusMomentoLogger) Trace(message string, args ...string) {
	l.logrusLogger.Tracef(message, args)
}

func (l LogrusMomentoLogger) Debug(message string, args ...string) {
	l.logrusLogger.Debugf(message, args)
}

func (l LogrusMomentoLogger) Info(message string, args ...string) {
	l.logrusLogger.Infof(message, args)
}

func (l LogrusMomentoLogger) Warn(message string, args ...string) {
	l.logrusLogger.Warnf(message, args)
}

func (l LogrusMomentoLogger) Error(message string, args ...string) {
	l.logrusLogger.Errorf(message, args)
}

type LogrusMomentoLoggerFactory struct{}

func NewLogrusMomentoLoggerFactory() logger.MomentoLoggerFactory {
	return &LogrusMomentoLoggerFactory{}
}

func (lf LogrusMomentoLoggerFactory) GetLogger(loggerName string) logger.MomentoLogger {
	log.SetFlags(0)
	return &LogrusMomentoLogger{logrusLogger: logrus.WithFields(logrus.Fields{"library": "Momento", "logger": loggerName})}
}
