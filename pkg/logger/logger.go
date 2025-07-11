package logger

import (
	"io"
	"log/slog"
	"time"

	"github.com/charmbracelet/log"
)

func NewHandler(logLevel slog.Level, writer io.Writer) *log.Logger {
	var lvl log.Level
	switch logLevel {
	case slog.LevelDebug:
		lvl = log.DebugLevel
	case slog.LevelInfo:
		lvl = log.InfoLevel
	case slog.LevelWarn:
		lvl = log.WarnLevel
	case slog.LevelError:
		lvl = log.ErrorLevel
	}

	logger := log.NewWithOptions(writer, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		CallerFormatter: log.ShortCallerFormatter,
		TimeFormat:      time.TimeOnly,
		Level:           lvl,
	})
	return logger
}
