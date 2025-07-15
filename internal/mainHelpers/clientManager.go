package mainhelpers

import (
	"log/slog"
	"os"

	filedownloader "github.com/volodymyrzuyev/goCsInspect/internal/fileDownloader"
	sqlite "github.com/volodymyrzuyev/goCsInspect/internal/storage/sqlite"
	"github.com/volodymyrzuyev/goCsInspect/pkg/clientmanagement"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
	"github.com/volodymyrzuyev/goCsInspect/pkg/gamefileupdater"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage"
)

func InitDefaultClientManager(
	cfg config.Config,
	mainLogger *slog.Logger,
	cleanLogger *slog.Logger,
) clientmanagement.ClientManager {

	storage, err := sqlite.NewSQLiteStore(cfg.DatabaseString, cleanLogger)
	if err != nil {
		mainLogger.Error("unable to connect to database, stoping", "error", err)
		os.Exit(1)
	}

	return InitClientManager(storage, cfg, mainLogger, cleanLogger)
}

func InitClientManagerNoStorage(
	cfg config.Config,
	mainLogger *slog.Logger,
	cleanLogger *slog.Logger,
) clientmanagement.ClientManager {

	return InitClientManager(&dummyStorage{}, cfg, mainLogger, cleanLogger)
}

func InitClientManager(
	str storage.Storage,
	cfg config.Config,
	mainLogger *slog.Logger,
	cleanLogger *slog.Logger,
) clientmanagement.ClientManager {

	fileDownloader := filedownloader.NewFileDownloader(cleanLogger)
	fileManager := gamefileupdater.NewFileUpdater(
		cfg.GameFilesAutoUpdateInverval,
		cfg.AutoUpdateGameFiles,
		cfg.GameLanguageLocation,
		cfg.GameItemsLocation,
		fileDownloader,
		cleanLogger,
	)
	gameItems, err := fileManager.UpdateFiles()
	if err != nil {
		mainLogger.Error("unable to generate new game items, stoping", "error", err)
		os.Exit(1)
	}
	det, err := detailer.NewDetailerWithCSItems(gameItems, cleanLogger)
	if err != nil {
		mainLogger.Error("unable to create new item detailer, stoping")
		os.Exit(1)
	}
	fileManager.RegisterDetailer(det)

	cm, err := clientmanagement.NewClientManager(
		cfg.RequestTTl,
		cfg.ClientCooldown,
		det,
		str,
		cleanLogger,
	)
	if err != nil {
		mainLogger.Error("unable to create new client manager, stoping", "error", err)
		os.Exit(1)
	}

	return cm
}
