package core

import "github.com/sirupsen/logrus"

var Log = initializeLogger()

func initializeLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	return log
}
