package internal

import "github.com/sirupsen/logrus"

var logger = logrus.New()

func init() {
	logger.Formatter = &logrus.JSONFormatter{}
}

// NewLogger returns configured logrus instanse.
func NewLogger() *logrus.Logger {
	return logger
}
