package cmd

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTraceLogger(t *testing.T) {
	tests := []struct {
		name        string
		level       log.Level
		callerLevel int
	}{
		{
			name:        "debug level trace logger",
			level:       log.DebugLevel,
			callerLevel: 2,
		},
		{
			name:        "info level trace logger",
			level:       log.InfoLevel,
			callerLevel: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTraceLogger(tt.level, tt.callerLevel)

			require.NotNil(t, l)
			assert.True(t, l.trace)
			assert.Equal(t, tt.callerLevel, l.callerLevel)
			assert.Equal(t, tt.level, l.Level)
		})
	}
}

func TestNewOutputLogger(t *testing.T) {
	tests := []struct {
		name        string
		callerLevel int
	}{
		{
			name:        "default output logger",
			callerLevel: 1,
		},
		{
			name:        "output logger with custom caller level",
			callerLevel: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewOutputLogger(tt.callerLevel)

			require.NotNil(t, l)
			assert.False(t, l.trace)
			assert.Equal(t, tt.callerLevel, l.callerLevel)
			assert.Equal(t, log.InfoLevel, l.Level)
		})
	}
}

func TestLoggerGenerateFields(t *testing.T) {
	tests := []struct {
		name       string
		logger     *Logger
		wantTrace  bool
		wantFields []string
	}{
		{
			name:       "trace logger includes file and line",
			logger:     NewTraceLogger(log.DebugLevel, 0),
			wantTrace:  true,
			wantFields: []string{"file", "line"},
		},
		{
			name:       "output logger has empty fields",
			logger:     NewOutputLogger(0),
			wantTrace:  false,
			wantFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.logger.generateFields()

			if tt.wantTrace {
				assert.Contains(t, fields, "file")
				assert.Contains(t, fields, "line")
			} else {
				assert.Empty(t, fields)
			}
		})
	}
}

func TestLoggerOutputMethods(t *testing.T) {
	tests := []struct {
		name     string
		logFunc  func(l *Logger)
		level    log.Level
		contains string
	}{
		{
			name:     "Debug message",
			logFunc:  func(l *Logger) { l.Debug("debug msg") },
			level:    log.DebugLevel,
			contains: "debug msg",
		},
		{
			name:     "Debugf message",
			logFunc:  func(l *Logger) { l.Debugf("debug %s", "formatted") },
			level:    log.DebugLevel,
			contains: "debug formatted",
		},
		{
			name:     "Info message",
			logFunc:  func(l *Logger) { l.Info("info msg") },
			level:    log.InfoLevel,
			contains: "info msg",
		},
		{
			name:     "Infof message",
			logFunc:  func(l *Logger) { l.Infof("info %s", "formatted") },
			level:    log.InfoLevel,
			contains: "info formatted",
		},
		{
			name:     "Error message",
			logFunc:  func(l *Logger) { l.Error("error msg") },
			level:    log.ErrorLevel,
			contains: "error msg",
		},
		{
			name:     "Errorf message",
			logFunc:  func(l *Logger) { l.Errorf("error %s", "formatted") },
			level:    log.ErrorLevel,
			contains: "error formatted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewOutputLogger(0)
			l.Level = log.DebugLevel

			var buf bytes.Buffer
			l.SetOutput(&buf)

			tt.logFunc(l)

			assert.Contains(t, buf.String(), tt.contains)
		})
	}
}

func TestTraceLoggerIncludesCallerInfo(t *testing.T) {
	l := NewTraceLogger(log.DebugLevel, 0)
	l.SetFormatter(&log.JSONFormatter{})

	var buf bytes.Buffer
	l.SetOutput(&buf)

	l.Debug("trace test")

	output := buf.String()
	assert.Contains(t, output, "file")
	assert.Contains(t, output, "line")
	assert.Contains(t, output, "trace test")
}

func TestOutputLoggerExcludesCallerInfo(t *testing.T) {
	l := NewOutputLogger(0)
	l.SetFormatter(&log.JSONFormatter{})

	var buf bytes.Buffer
	l.SetOutput(&buf)

	l.Info("output test")

	output := buf.String()
	assert.NotContains(t, output, "\"file\"")
	assert.Contains(t, output, "output test")
}
