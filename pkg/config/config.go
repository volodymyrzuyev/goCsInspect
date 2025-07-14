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

	LogLevel slog.Level

	BindIP string
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

		LogLevel: slog.LevelDebug,

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
	cfg.GameLanguageLocation = common.GetAbsolutePath(cfg.GameLanguageLocation)
	cfg.DatabaseString = common.GetAbsolutePath(cfg.DatabaseString)

	return cfg, nil
}
