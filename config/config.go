package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/volodymyrzuyev/goCsInspect/pkg/common"
)

var (
	TimeOutDuration = 10 * time.Second
	RequestCooldown = 2 * time.Second

	// Use relative paths from project root
	EnglishFile = "game_files/csgo_english.txt"
	GameItems   = "game_files/items_game.txt"

	IsDebug       = true
	DebugLocation = "debug/logs"
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
	// Ensure debug directory exists
	debugDir := filepath.Dir(common.GetAbsolutePath(DebugLocation))
	if err := os.MkdirAll(debugDir, 0755); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(common.GetAbsolutePath(DebugLocation), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	return logger
}

// GetEnglishFile returns the absolute path to the English file
func GetEnglishFile() string {
	return common.GetAbsolutePath(EnglishFile)
}

// GetGameItems returns the absolute path to the game items file
func GetGameItems() string {
	return common.GetAbsolutePath(GameItems)
}
