package common

import (
	"path/filepath"
	"runtime"
)


func GetProjectRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current file path")
	}

	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	return projectRoot
}

func GetAbsolutePath(relativePath string) string {
	projectRoot := GetProjectRoot()
	return filepath.Join(projectRoot, relativePath)
}
