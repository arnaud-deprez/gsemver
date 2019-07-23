package log

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
	f := &log.TextFormatter{}
	f.ForceColors = true
	log.SetFormatter(f)
}

// Level type
type Level uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// FatalLevel level, highest level of severity. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel Level = iota + 1
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Trace log formatted message at trace level
func Trace(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

// Debug log formatted message at debug level
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info log formatted message at info level
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn log formatted message at info level
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error log formatted message at error level
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal log formatted message at fatal level
// Calls os.Exit(1) after logging
func Fatal(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	log.SetLevel(log.Level(level))
}

// SetLevelS sets the standard logger level from string.
// Level can be: trace, debug, info, error or fatal
func SetLevelS(level string) {
	l, err := log.ParseLevel(strings.ToLower(level))
	if err != nil {
		Fatal("Cannot configure logger caused by %v", err)
	}
	log.SetLevel(l)
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level Level) bool {
	return log.IsLevelEnabled(log.Level(level))
}
