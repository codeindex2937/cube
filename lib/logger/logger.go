package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Log() *logrus.Logger {
	if log == nil {
		log = logrus.New()

		log.Out = os.Stdout
		log.SetLevel(logrus.DebugLevel)
		log.SetFormatter(&logrus.TextFormatter{})
	}
	return log
}
