package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/volodymyrzuyev/goCsInspect/pkg/common"
	"github.com/volodymyrzuyev/goCsInspect/pkg/creds"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RequestTTl     time.Duration
	ClientCooldown time.Duration
	Accounts       []creds.Credentials

	GameItemsLocation           string
	GameLanguageLocation        string
	AutoUpdateGameFiles         bool
	GameFilesAutoUpdateInverval time.Duration

	DatabaseString string

	LogLevel string

	BindIP string
}

func (cfg *Config) VertifyConfig() error {
	var err error

	cfg.GameItemsLocation = common.GetAbsolutePath(cfg.GameItemsLocation)
	if err = common.VertifyAndCreateFile(cfg.GameItemsLocation); err != nil {
		return nil
	}

	cfg.GameLanguageLocation = common.GetAbsolutePath(cfg.GameLanguageLocation)
	if err = common.VertifyAndCreateFile(cfg.GameLanguageLocation); err != nil {
		return nil
	}

	cfg.DatabaseString = common.GetAbsolutePath(cfg.DatabaseString)
	if err = common.VertifyAndCreateFile(cfg.DatabaseString); err != nil {
		return nil
	}

	return nil
}

func (cfg *Config) GetLogLevel() slog.Level {
	switch cfg.LogLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

var (
	DefaultConfig = Config{
		RequestTTl:     3 * time.Second,
		ClientCooldown: 1*time.Second + 100*time.Millisecond,
		Accounts:       []creds.Credentials{},

		GameItemsLocation:           common.GetAbsolutePath("game_files/items_game.txt"),
		GameLanguageLocation:        common.GetAbsolutePath("game_files/csgo_english.txt"),
		AutoUpdateGameFiles:         true,
		GameFilesAutoUpdateInverval: 4 * time.Hour,

		DatabaseString: common.GetAbsolutePath("data.db"),

		LogLevel: "INFO",

		BindIP: "0.0.0.0:8080",
	}

	DefaultConfigLocation = common.GetAbsolutePath("config.yaml")
)

func ParseConfig(relativePath string) (Config, error) {
	out, err := os.ReadFile(relativePath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = yaml.Unmarshal(out, &cfg)
	if err != nil {
		return Config{}, err
	}

	cfg.GameItemsLocation = common.GetAbsolutePath(cfg.GameItemsLocation)
	if err = common.VertifyAndCreateFile(cfg.GameItemsLocation); err != nil {
		return Config{}, nil
	}

	cfg.GameLanguageLocation = common.GetAbsolutePath(cfg.GameLanguageLocation)
	if err = common.VertifyAndCreateFile(cfg.GameLanguageLocation); err != nil {
		return Config{}, nil
	}

	cfg.DatabaseString = common.GetAbsolutePath(cfg.DatabaseString)
	if err = common.VertifyAndCreateFile(cfg.DatabaseString); err != nil {
		return Config{}, nil
	}

	return cfg, nil
}
