package log

import (
	log "github.com/sirupsen/logrus"
)

//Debug log formatted message at debug level
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

//Info log formatted message at info level
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

//Warn log formatted message at info level
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

//Error log formatted message at error level
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}