package logger

import (
	"io"
	"log/slog"
	"time"

	"github.com/charmbracelet/log"
)

type Logger interface {
	Debug(msg any, keyvals ...any)
	Info(msg any, keyvals ...any)
	Warn(msg any, keyvals ...any)
	Error(msg any, keyvals ...any)
	SetPrefix(string)
}

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
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           lvl,
	})
	return logger
}
