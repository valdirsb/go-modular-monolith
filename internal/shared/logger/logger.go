package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type logger struct {
	*log.Logger
}

func NewLogger() Logger {
	return &logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *logger) Info(msg string) {
	l.Logger.Println("INFO: " + msg)
}

func (l *logger) Error(msg string) {
	l.Logger.Println("ERROR: " + msg)
}