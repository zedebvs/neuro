package utils

import (
	"os"
	"path/filepath"
)

var currentDir string

func GetPath(parts []string) (string, error) {
	PartsSlice := append([]string{currentDir}, parts...)
	return filepath.Join(PartsSlice...), nil
}

func CreateDir(path string) error {
	FileDir := filepath.Dir(path)
	return os.MkdirAll(FileDir, 0756)
}

func init() {
	cur, err := os.Executable()
	if err != nil {
		panic("Не удалось найти путь к main")
	}
	currentDir = filepath.Dir(cur)
}
