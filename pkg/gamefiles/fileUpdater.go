package gamefiles

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/volodymyrzuyev/go-csgo-item-parser/csgo"
	"github.com/volodymyrzuyev/go-csgo-item-parser/parser"
	det "github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
)

type Downloader interface {
	GetFiles() (FilePaths, error)
}

type Updater interface {
	UpdateFiles() (*csgo.Csgo, error)
	RegisterDetailer(det.Detailer)
}

type FilePaths struct {
	LanguageFile string
	GameItems    string
}

type updater struct {
	autoUpdate     bool
	updateInterval time.Duration
	languagePath   string
	gameItemsPath  string

	d det.Detailer
	f Downloader
	l *slog.Logger
}

func NewUpdater(
	updateInterval time.Duration,
	autoUpdate bool,
	languageFilePath, gameItemsPath string,
	downloader Downloader,
) Updater {

	fu := &updater{
		updateInterval: updateInterval,
		autoUpdate:     autoUpdate,
		languagePath:   languageFilePath,
		gameItemsPath:  gameItemsPath,

		f: downloader,
		l: slog.Default().WithGroup("FileUpdater"),
	}

	return fu
}

func (u *updater) UpdateFiles() (*csgo.Csgo, error) {
	var newItems *csgo.Csgo
	var err error

	if u.autoUpdate {
		fp, err := u.f.GetFiles()
		if err != nil {
			u.l.Error("could not get new files", "error", err)
			return nil, err
		}

		newItems, err = u.getCsItems(fp)
		if err != nil {
			return nil, err
		}

		err = replaceFiles(fp.LanguageFile, u.languagePath)
		if err != nil {
			u.l.Error("could not replace language file", "error", err)
			return nil, err
		}

		err = replaceFiles(fp.GameItems, u.gameItemsPath)
		if err != nil {
			u.l.Error("could not replace game items file", "error", err)
			return nil, err
		}

		u.l.Info("game files updated")
	} else {
		newItems, err = u.getCsItems(FilePaths{LanguageFile: u.languagePath, GameItems: u.gameItemsPath})
		if err != nil {
			return nil, err
		}
		u.l.Info("original files parsed")
	}

	return newItems, nil
}

func (u *updater) RegisterDetailer(d det.Detailer) {
	u.d = d
	go u.autoLoop()
}

func (u *updater) getCsItems(fp FilePaths) (*csgo.Csgo, error) {
	languageData, err := parser.Parse(fp.LanguageFile)
	if err != nil {
		u.l.Error("could not parser language file", "error", err)
		return nil, err
	}

	itemData, err := parser.Parse(fp.GameItems)
	if err != nil {
		u.l.Error("could not parser item file", "error", err)
		return nil, err
	}

	allItems, err := csgo.New(languageData, itemData)
	if err != nil {
		u.l.Error("could not parser cs files", "error", err)
		return nil, err
	}

	return allItems, nil
}

func (u *updater) autoLoop() {
	for {
		time.Sleep(u.updateInterval)
		newItems, err := u.UpdateFiles()
		if err != nil {
			u.l.Error("could not update files")
			continue
		}

		u.d.UpdateItems(newItems)
	}
}

func replaceFiles(sourcePath, destPath string) error {
	err := os.Rename(sourcePath, destPath)
	if err == nil {
		return nil
	}

	_, ok := err.(*os.LinkError)
	if !ok {
		return err
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		os.Remove(destPath)
		return err
	}

	if err = destFile.Sync(); err != nil {
		return err
	}

	if err = os.Remove(sourcePath); err != nil {
		return err
	}

	return nil
}
