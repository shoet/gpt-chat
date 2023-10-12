package util

import (
	"os"
	"path/filepath"
)

func GetProjectRoot(currentDir string) (string, error) {
	files, err := os.ReadDir(currentDir)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if f.Name() == "go.mod" {
			return currentDir, nil
		}
	}
	upper := filepath.Dir(currentDir)
	return GetProjectRoot(upper)
}
