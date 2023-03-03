package logger

type NoopMomentoLogger struct {
}

func (*NoopMomentoLogger) Info(message string, args ...string) {
	// no-op
}

func (*NoopMomentoLogger) Debug(message string, args ...string) {
	// no-op
}

func (*NoopMomentoLogger) Warn(message string, args ...string) {
	// no-op
}

func (*NoopMomentoLogger) Error(message string, args ...string) {
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
