package tacks

import "github.com/Sirupsen/logrus"

var logger = logrus.New()

// Logger returns the tacks Logger instance
func Logger() *logrus.Logger {
	return logger
}
