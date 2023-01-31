package config

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(message string)
	Debug(message string)
}

type LoggerConfiguration struct {
	Name       string
}

type BuiltInLoggerClient struct {
	logger *log.Logger
	name   string
}

func NewMomentoLogger(l *LoggerConfiguration) Logger {
	if l.LoggerType == builtin {
		return &BuiltInLoggerClient{
			logger: l.,
			name:   l.Name,
		}
	}
	return &BuiltInLoggerClient{
		logger: log.Default(),
		name:   l.Name,
	}
}

func (lc *BuiltInLoggerClient) Info(message string) {
	str := fmt.Sprintf(`{"level": "INFO", "%s": "%s"}`, lc.name, message)
	log.SetFlags(0)
	log.Print(str)
}

func (lc *BuiltInLoggerClient) Debug(message string) {
	str := fmt.Sprintf(`{"level": "DEBUG", "%s": "%s"}`, lc.name, message)
	log.SetFlags(0)
	log.Print(str)
}
