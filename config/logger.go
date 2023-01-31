package config

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(message string)
	Debug(message string)
}

type logLevel string

const (
	info  logLevel = "INFO"
	debug logLevel = "DEBUG"
)

type logFormat string

const (
	json logFormat = "JSON"
)

type LoggerOptions struct {
	Level  logLevel
	Format logFormat
	Name   string
}

type BuiltInLoggerClient struct {
	logger    *log.Logger
	name      string
	logLevel  logLevel
	logFormat logFormat
}

type loggerType string

const (
	builtin loggerType = "BUILTIN"
)

func NewBuiltInLogger(l *LoggerOptions) Logger {
	return &BuiltInLoggerClient{
		logger:    log.Default(),
		name:      l.Name,
		logLevel:  l.Level,
		logFormat: l.Format,
	}
}

func (lc *BuiltInLoggerClient) Info(message string) {
	str := fmt.Sprintf(`{"level": "INFO", "message": "%s", "name": "%s""}`, message, lc.name)
	log.SetFlags(0)
	log.Print(str)
}

func (lc *BuiltInLoggerClient) Debug(message string) {
	str := fmt.Sprintf(`{"level": "DEBUG", "message": "%s", "name": "%s""}`, message, lc.name)
	log.SetFlags(0)
	log.Print(str)
}
