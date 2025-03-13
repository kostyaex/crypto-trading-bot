package utils

import (
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(level string) *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	logger.SetLevel(logLevel)

	return &Logger{Logger: logger}
}

func (l *Logger) Writer() *os.File {
	return os.Stdout
}
