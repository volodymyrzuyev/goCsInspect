package config

import (
	"log/slog"
	"os"
	"time"
)

var (
	TimeOutDuration = 10 * time.Second
	RequestCooldown = 10 * time.Second

	EnglishFile = "./game_files/csgo_english.txt"
	GameItems   = "./game_files/items_game.txt"

	IsDebug       = true
	DebugLocation = "./debug/logs"
	DebugLogger   = getDebugLogger()

	DefaultClientConfig = ClientConfig{
		RequestCooldown: RequestCooldown,
		TimeOutDuration: TimeOutDuration,
		IsDebug:         IsDebug,
		DebugLocation:   DebugLocation,
		DebugLogger:     DebugLogger,
	}
)

type ClientConfig struct {
	RequestCooldown time.Duration
	TimeOutDuration time.Duration

	IsDebug       bool
	DebugLocation string
	DebugLogger   *slog.Logger
}

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
