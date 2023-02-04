package logger

type NoopMomentoLogger struct {
}

func (*NoopMomentoLogger) Info(message string, args ...any) {
	// no-op
}

func (*NoopMomentoLogger) Debug(message string, args ...any) {
	// no-op
}

type NoopMomentoLoggerFactory struct {
}

func NewNoopMomentoLoggerFactory() MomentoLoggerFactory {
	return &NoopMomentoLoggerFactory{}
}

func (*NoopMomentoLoggerFactory) GetLogger(loggerName string) MomentoLogger {
	return &NoopMomentoLogger{}
}