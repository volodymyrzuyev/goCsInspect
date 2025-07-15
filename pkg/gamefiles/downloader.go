package gamefiles

import (
	"context"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
)

const (
	owner             = "csfloat"
	repo              = "cs-files"
	gameItemsRepoPath = "static/items_game.txt"
	englishRepoPath   = "static/csgo_english.txt"
)

type fileDownloader struct {
	client *github.Client

	gameItemsLastSHA string
	languageLastSHA  string

	l *slog.Logger
}

func NewFileDownloader() Downloader {

	return &fileDownloader{
		client: github.NewClient(nil),

		l: slog.Default().WithGroup("FileDownloader"),
	}
}

func (f *fileDownloader) GetFiles() (FilePaths, error) {
	gameItems, err := f.querryGithub(gameItemsRepoPath)
	if err != nil {
		f.l.Error("could not get new game items", "error", err)
		return FilePaths{}, err
	}

	language, err := f.querryGithub(englishRepoPath)
	if err != nil {
		f.l.Error("could not get new language file", "error", err)
		return FilePaths{}, err
	}

	if !f.compareShas(*language.SHA, *gameItems.SHA) {
		return FilePaths{}, errors.ErrNoNewFiles
	}

	languageData, err := downloadFileData(language)
	if err != nil {
		f.l.Error("could not get new language file", "error", err)
		return FilePaths{}, err
	}

	itemsData, err := downloadFileData(gameItems)
	if err != nil {
		f.l.Error("could not get new game items", "error", err)
		return FilePaths{}, err
	}

	languageFileName, err := createTemp(languageData)
	if err != nil {
		f.l.Error("could not save language file", "error", err)
		return FilePaths{}, err
	}

	itemsFileName, err := createTemp(itemsData)
	if err != nil {
		f.l.Error("could not save game items file", "error", err)
		return FilePaths{}, err
	}

	f.l.Info("new files downloaded")

	return FilePaths{LanguageFile: languageFileName, GameItems: itemsFileName}, nil
}

func (f *fileDownloader) querryGithub(repoPath string) (*github.RepositoryContent, error) {
	fileContent, _, _, err := f.client.Repositories.GetContents(
		context.Background(),
		owner,
		repo,
		repoPath,
		nil,
	)
	return fileContent, err
}

func (f *fileDownloader) compareShas(englishSHA, itemsSHA string) bool {
	if f.gameItemsLastSHA == itemsSHA && f.languageLastSHA == englishSHA {
		return false
	}

	return true
}

func downloadFileData(file *github.RepositoryContent) ([]byte, error) {
	if file.DownloadURL != nil && *file.DownloadURL != "" {
		resp, err := http.Get(*file.DownloadURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, err
		}

		return io.ReadAll(resp.Body)
	}

	return base64.StdEncoding.DecodeString(*file.Content)
}

func createTemp(data []byte) (string, error) {
	temp, err := os.CreateTemp("", "goCsInspect")
	if err != nil {
		return "", err
	}

	_, err = temp.Write(data)
	if err != nil {
		return "", err
	}

	err = temp.Close()
	if err != nil {
		return "", err
	}

	return temp.Name(), nil
}
