package app

import (
	"fmt"
	"io"
	"os"
	"time"

	"path/filepath"

	"github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

const timestampFormat = time.RFC3339

type LogConfig struct {
	// LogLevel represents application logging level.
	Level level `envconfig:"LOG_LEVEL" default:"INFO"`
	// LogHumanReadable indicates if logger should use coloured output
	// in logfmt format instead of default JSON.
	Human bool `envconfig:"LOG_HUMAN" default:"false"`
	// LogFile is os path to file where logs are printed. If empty logs
	// will are printed to standard output.
	File string `envconfig:"LOG_FILE"`
}

type level struct {
	logrus.Level
}

// Decode decodes log level
func (l *level) Decode(s string) error {
	var err error
	l.Level, err = logrus.ParseLevel(s)
	return err
}

func newLog(cfg LogConfig, name string) *logrus.Entry {
	l := logrus.New()
	l.Out = logWriter(cfg.File)
	l.Level = cfg.Level.Level
	if cfg.Human {
		return logrus.NewEntry(l)
	}
	// set formatter for json output
	l.Formatter = &logrus.JSONFormatter{
		TimestampFormat: timestampFormat,
		FieldMap:        logrus.FieldMap{logrus.FieldKeyTime: "ts"},
	}
	// add path and stack logging
	l.Hooks.Add(logrus_stack.StandardHook())
	// add service name field
	return l.WithField("service", name)
}

func logWriter(path string) io.Writer {
	if path == "" {
		return os.Stdout
	}
	// create directory if not exists
	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(filepath.Dir(path), 0776)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file %s: %s", path, err))
	}
	return file
}
