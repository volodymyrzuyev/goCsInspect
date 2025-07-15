package common

import (
	"os"
	"path/filepath"
)

func GetAbsolutePath(relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	projectRoot, err := os.Getwd()
	if err != nil {
		panic("could not get working directory")
	}

	return filepath.Join(projectRoot, relativePath)
}

func CreateFile(path string) error {
	fullPath := GetAbsolutePath(path)

	err := os.MkdirAll(filepath.Dir(path), 0770)
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	return file.Close()
}

func VertifyAndCreateFile(path string) error {
	_, err := os.Stat(GetAbsolutePath(path))
	if os.IsNotExist(err) {
		return CreateFile(path)
	}

	return err
}
