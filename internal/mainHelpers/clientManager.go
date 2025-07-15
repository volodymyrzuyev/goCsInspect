package mainhelpers

import (
	"log/slog"
	"os"

	"github.com/volodymyrzuyev/goCsInspect/pkg/clientmanagement"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
	"github.com/volodymyrzuyev/goCsInspect/pkg/gamefiles"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage/sqlite"
)

func InitDefaultClientManager(cfg config.Config) clientmanagement.Manager {
	mainLogger := slog.Default().WithGroup("Main")

	storage, err := sqlite.NewSQLiteStore(cfg.DatabaseString)
	if err != nil {
		mainLogger.Error("unable to connect to database, stoping", "error", err)
		os.Exit(1)
	}

	return InitClientManager(storage, cfg)
}

func InitClientManagerNoStorage(cfg config.Config) clientmanagement.Manager {

	return InitClientManager(&dummyStorage{}, cfg)
}

func InitClientManager(str storage.Storage, cfg config.Config) clientmanagement.Manager {
	mainLogger := slog.Default().WithGroup("Main")

	downloader := gamefiles.NewFileDownloader()
	updater := gamefiles.NewUpdater(
		cfg.GameFilesAutoUpdateInverval,
		cfg.AutoUpdateGameFiles,
		cfg.GameLanguageLocation,
		cfg.GameItemsLocation,
		downloader,
	)
	gameItems, err := updater.UpdateFiles()
	if err != nil {
		mainLogger.Error("unable to generate new game items, stoping", "error", err)
		os.Exit(1)
	}
	det, err := detailer.NewDetailerWithCSItems(gameItems)
	if err != nil {
		mainLogger.Error("unable to create new item detailer, stoping")
		os.Exit(1)
	}
	updater.RegisterDetailer(det)

	cm, err := clientmanagement.New(
		cfg.RequestTTl,
		cfg.ClientCooldown,
		det,
		str,
	)
	if err != nil {
		mainLogger.Error("unable to create new client manager, stoping", "error", err)
		os.Exit(1)
	}

	return cm
}
