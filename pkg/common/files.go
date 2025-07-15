package common

import (
	"os"
	"path/filepath"
)

// Returns absolute path of a file/directory. If the relativePath is an absolute
// path, returns the path passed in. 
// Panics if unable to get PWD
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

// Creates a file with any directory's the path outlines
// Returns any error that is encountered
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

// Checks if a path exists, if not, creates it
func VertifyAndCreateFile(path string) error {
	_, err := os.Stat(GetAbsolutePath(path))
	if os.IsNotExist(err) {
		return CreateFile(path)
	}

	return err
}
