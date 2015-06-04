package tacks

import "github.com/Sirupsen/logrus"

var logger *logrus.Logger = logrus.New()

const (
	Nothing = ""
)

func Logger() *logrus.Logger {
	return logger
}
