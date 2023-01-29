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
	Name string
}

type LoggerClient struct {
	logger *log.Logger
	name   string
}

func NewLogger(l *LoggerConfiguration) Logger {
	return &LoggerClient{
		logger: log.Default(),
		name:   l.Name,
	}
}

func (lc *LoggerClient) Info(message string) {
	str := fmt.Sprintf(`{"level": "INFO", "%s": "%s"}`, lc.name, message)
	log.SetFlags(0)
	log.Print(str)
}

func (lc *LoggerClient) Debug(message string) {
	str := fmt.Sprintf(`{"level": "DEBUG", "%s": "%s"}`, lc.name, message)
	log.SetFlags(0)
	log.Print(str)
}
