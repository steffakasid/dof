package cmd

import (
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type Logger struct {
	*log.Logger
	trace       bool
	callerLevel int
}

func NewTraceLogger(loglevel log.Level, callerLevel int) *Logger {
	logger := &Logger{log.New(), true, callerLevel}
	logger.Level = loglevel

	return logger
}

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

func (l *Logger) Debug(arg ...interface{}) {
	l.WithFields(l.generateFields()).Debug(arg...)
}

func (l *Logger) Debugf(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Debugf(format, arg...)
}

func (l *Logger) Info(arg ...interface{}) {
	l.WithFields(l.generateFields()).Info(arg...)
}

func (l *Logger) Infof(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Infof(format, arg...)
}

func (l *Logger) Error(arg ...interface{}) {
	l.WithFields(l.generateFields()).Error(arg...)
}

func (l *Logger) Errorf(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Errorf(format, arg...)
}

func (l *Logger) Fatal(arg ...interface{}) {
	l.WithFields(l.generateFields()).Fatal(arg...)
}

func (l *Logger) Fatalf(format string, arg ...interface{}) {
	l.WithFields(l.generateFields()).Fatalf(format, arg...)
}
