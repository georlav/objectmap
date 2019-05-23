// Package cli outputs messages to stdout
package cli

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewOutput creates and returns a new logger instance that logs data to std out
func NewOutput() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.FatalLevel)

	f := logrus.TextFormatter{}
	f.DisableLevelTruncation = false
	f.ForceColors = true
	f.FullTimestamp = false
	f.DisableTimestamp = true
	f.EnvironmentOverrideColors = false

	logger.SetFormatter(&f)

	return logger
}
