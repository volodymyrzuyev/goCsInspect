package logger

import (
	"fmt"
	"io"
	"log"
)

type Logger interface {
	Error(string, ...any)
	Info(string, ...any)
	Debug(string, ...any)
}

type logger struct {
	error *log.Logger
	info  *log.Logger
	debug *log.Logger
}

// If error logger is provided, function will panic
func (l *logger) Error(format string, v ...any) {
	l.write(l.error, format, v...)
}

// If info logger is provided, function will panic
func (l *logger) Info(format string, v ...any) {
	l.write(l.info, format, v...)
}

// If debug logger is provided, function will panic
func (l *logger) Debug(format string, v ...any) {
	l.write(l.debug, format, v...)
}

func (l *logger) write(writer *log.Logger, format string, v ...any) {
	if writer == nil {
		panic("You need to initilize the loggers with logger.NewLogger() in order to use it. Some of the loggers are not initialized")
	}
	writer.Output(3, fmt.Sprintf(format, v...))
}

// Creates a logger that will output to the provided io.Writer
func NewLogger(output io.Writer) Logger {
	log := logger{
		error: log.New(output, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		info:  log.New(output, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		debug: log.New(output, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return &log
}
