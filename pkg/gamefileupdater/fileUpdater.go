package gamefileupdater

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/volodymyrzuyev/go-csgo-item-parser/csgo"
	"github.com/volodymyrzuyev/go-csgo-item-parser/parser"
	det "github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
)

type FileDownloader interface {
	GetFiles() (FilePaths, error)
}

type FileUpdater interface {
	UpdateFiles() (*csgo.Csgo, error)
	RegisterDetailer(det.Detailer)
}

type FilePaths struct {
	LanguageFile string
	GameItems    string
}

type fileUpdater struct {
	autoUpdate     bool
	updateInterval time.Duration
	languagePath   string
	gameItemsPath  string

	d det.Detailer
	f FileDownloader
	l *slog.Logger
}

func NewFileUpdater(
	updateInterval time.Duration,
	autoUpdate bool,
	languageFilePath, gameItemsPath string,
	fileDownloader FileDownloader,
	l *slog.Logger,
) FileUpdater {

	fu := &fileUpdater{
		updateInterval: updateInterval,
		autoUpdate:     autoUpdate,
		languagePath:   languageFilePath,
		gameItemsPath:  gameItemsPath,

		f: fileDownloader,
		l: l.WithGroup("FileUpdater"),
	}

	return fu
}

func (f *fileUpdater) UpdateFiles() (*csgo.Csgo, error) {
	var newItems *csgo.Csgo
	var err error

	if f.autoUpdate {
		fp, err := f.f.GetFiles()
		if err != nil {
			f.l.Error("could not get new files", "error", err)
			return nil, err
		}

		newItems, err = f.getCsItems(fp)
		if err != nil {
			return nil, err
		}

		err = replaceFiles(fp.LanguageFile, f.languagePath)
		if err != nil {
			f.l.Error("could not replace language file", "error", err)
			return nil, err
		}

		err = replaceFiles(fp.GameItems, f.gameItemsPath)
		if err != nil {
			f.l.Error("could not replace game items file", "error", err)
			return nil, err
		}

		f.l.Info("game files updated")
	} else {
		newItems, err = f.getCsItems(FilePaths{LanguageFile: f.languagePath, GameItems: f.gameItemsPath})
		if err != nil {
			return nil, err
		}
		f.l.Info("original files parsed")
	}

	return newItems, nil
}

func (f *fileUpdater) RegisterDetailer(d det.Detailer) {
	f.d = d
	go f.autoLoop()
}

func (f *fileUpdater) getCsItems(fp FilePaths) (*csgo.Csgo, error) {
	languageData, err := parser.Parse(fp.LanguageFile)
	if err != nil {
		f.l.Error("could not parser language file", "error", err)
		return nil, err
	}

	itemData, err := parser.Parse(fp.GameItems)
	if err != nil {
		f.l.Error("could not parser item file", "error", err)
		return nil, err
	}

	allItems, err := csgo.New(languageData, itemData)
	if err != nil {
		f.l.Error("could not parser cs files", "error", err)
		return nil, err
	}

	return allItems, nil
}

func (f *fileUpdater) autoLoop() {
	for {
		time.Sleep(f.updateInterval)
		newItems, err := f.UpdateFiles()
		if err != nil {
			f.l.Error("could not update files")
			continue
		}

		f.d.UpdateItems(newItems)
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
