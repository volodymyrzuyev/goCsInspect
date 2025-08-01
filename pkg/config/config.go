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
	// Amount of time for a request to live, if exceeded, request get's aborted
	RequestTTl time.Duration
	// Time between requests can be sent from a single client
	ClientCooldown time.Duration
	// Client credentials
	Accounts []creds.Account

	// Location of items_game.txt
	GameItemsLocation string
	// Location of language_(name).txt
	GameLanguageLocation string
	// Flag to enable/disable auto file updates
	AutoUpdateGameFiles bool
	// Interval of fetching new updates
	GameFilesAutoUpdateInverval time.Duration

	// Connection path for a database
	DatabaseString string

	// Log lever
	// [DEBUG, WARN, INFO, ERROR]
	LogLevel string

	// IP to which HTTP api will bind
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
		Accounts:       []creds.Account{},

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

// Parses config form a YAML file
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
