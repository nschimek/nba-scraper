package context

import "github.com/sirupsen/logrus"

func createLogger() *logrus.Logger {
	log := logrus.New()

	return log
}
