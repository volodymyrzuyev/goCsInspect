package config

import (
	"log/slog"
	"os"
	"time"
)

var (
	TimeOutDuration = 10 * time.Second
	RequestCooldown = 10 * time.Second

	IsDebug       = true
	DebugLocation = "./debug/logs"
	DebugLogger   = getDebugLogger()
)

func getDebugLogger() *slog.Logger {
	logFile, err := os.OpenFile(DebugLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	return logger
}
