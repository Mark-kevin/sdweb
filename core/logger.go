package core

import "github.com/sirupsen/logrus"

func LogInfoInit() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)            // 设置日志级别
	log.SetFormatter(&logrus.JSONFormatter{}) // 设置日志格式
	return log
}

func LogDebugInit() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)           // 设置日志级别
	log.SetFormatter(&logrus.JSONFormatter{}) // 设置日志格式
	return log
}

func LogErrorInit() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)           // 设置日志级别
	log.SetFormatter(&logrus.JSONFormatter{}) // 设置日志格式
	return log
}
