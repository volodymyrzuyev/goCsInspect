package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	filedownloader "github.com/volodymyrzuyev/goCsInspect/internal/fileDownloader"
	storage "github.com/volodymyrzuyev/goCsInspect/internal/storage/sqlite"
	"github.com/volodymyrzuyev/goCsInspect/internal/web"
	"github.com/volodymyrzuyev/goCsInspect/pkg/clientmanagement"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
	"github.com/volodymyrzuyev/goCsInspect/pkg/gamefileupdater"
	"github.com/volodymyrzuyev/goCsInspect/pkg/logger"
)

var (
	cfgLocation string
)

func init() {
	flag.StringVar(
		&cfgLocation, "config", config.DefaultConfigLocation, "configuration file used for the api",
	)
}

func main() {
	flag.Parse()

	cfg, err := config.ParseConfig(cfgLocation)
	if err != nil {
		fmt.Println("invalid configuration location, stoping")
		os.Exit(1)
	}

	l := slog.New(logger.NewHandler(cfg.LogLevel, os.Stdout))
	slog.SetDefault(l)
	lt := l.WithGroup("Main")

	storage, err := storage.NewSQLiteStore(cfg.DatabaseString, l)
	if err != nil {
		lt.Error("unable to connect to database, stoping", "error", err)
		os.Exit(1)
	}

	fileDownloader := filedownloader.NewFileDownloader(l)
	fileManager := gamefileupdater.NewFileUpdater(
		cfg.GameFilesAutoUpdateInverval,
		cfg.AutoUpdateGameFiles,
		cfg.GameLanguageLocation,
		cfg.GameItemsLocation,
		fileDownloader,
		l,
	)
	gameItems, err := fileManager.UpdateFiles()
	if err != nil {
		lt.Error("unable to generate new game items, stoping", "error", err)
		os.Exit(1)
	}
	det, err := detailer.NewDetailerWithCSItems(gameItems, l)
	if err != nil {
		lt.Error("unable to create new item detailer, stoping")
		os.Exit(1)
	}
	fileManager.RegisterDetailer(det)

	cm, err := clientmanagement.NewClientManager(
		cfg.RequestTTl,
		cfg.ClientCooldown,
		det,
		storage,
		l,
	)
	if err != nil {
		lt.Error("unable to create new client manager, stoping", "error", err)
		os.Exit(1)
	}

	for _, cli := range cfg.Accounts {
		err := cm.AddClient(cli)
		if err != nil {
			lt.Warn(
				fmt.Sprintf("client %v unable to login, won't be used", cli.Username),
				"error",
				err,
			)
		}
	}

	server := web.NewServer(cm, l)

	server.Run(cfg.BindIP)
}
