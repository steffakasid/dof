package cmd

import (
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with optional trace-level caller information.
type Logger struct {
	*log.Logger
	trace       bool
	callerLevel int
}

// NewTraceLogger creates a Logger that includes file/line info in log output.
func NewTraceLogger(loglevel log.Level, callerLevel int) *Logger {
	logger := &Logger{log.New(), true, callerLevel}
	logger.Level = loglevel

	return logger
}

// NewOutputLogger creates a Logger for standard output without trace info.
func NewOutputLogger(callerLevel int) *Logger {
	logger := &Logger{log.New(), false, callerLevel}
	logger.Level = log.InfoLevel

	return logger
}

func (l *Logger) generateFields() log.Fields {
	_, file, line, _ := runtime.Caller(l.callerLevel + 1)
	if l.trace {
		return log.Fields{"line": line, "file": path.Base(file)}
	}
	return log.Fields{}
}

// Debug logs a debug message with optional caller fields.
func (l *Logger) Debug(arg ...interface{}) {
	l.WithFields(l.generateFields()).Debug(arg...)
}

// Debugf logs a formatted debug message with optional caller fields.
func (l *Logger) Debugf(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Debugf(format, arg...)
}

// Info logs an info message with optional caller fields.
func (l *Logger) Info(arg ...interface{}) {
	l.WithFields(l.generateFields()).Info(arg...)
}

// Infof logs a formatted info message with optional caller fields.
func (l *Logger) Infof(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Infof(format, arg...)
}

// Error logs an error message with optional caller fields.
func (l *Logger) Error(arg ...interface{}) {
	l.WithFields(l.generateFields()).Error(arg...)
}

// Errorf logs a formatted error message with optional caller fields.
func (l *Logger) Errorf(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Errorf(format, arg...)
}

// Fatal logs a fatal message with optional caller fields and exits.
func (l *Logger) Fatal(arg ...interface{}) {
	l.WithFields(l.generateFields()).Fatal(arg...)
}

// Fatalf logs a formatted fatal message with optional caller fields and exits.
func (l *Logger) Fatalf(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Fatalf(format, arg...)
}
