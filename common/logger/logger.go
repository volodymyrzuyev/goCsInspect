package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func NewHandler() *log.Logger {
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		CallerFormatter: log.ShortCallerFormatter,
		TimeFormat:      time.TimeOnly,
		Level:           log.DebugLevel,
	})
	return logger
}
