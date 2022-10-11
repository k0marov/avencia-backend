package static_store

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// RecursiveDirCreator os.MkdirAll implements this
type RecursiveDirCreator = func(path string, perm fs.FileMode) error

// FileCreator os.WriteFile implements this
type FileCreator = func(name string, data []byte, perm fs.FileMode) error

func NewStaticFileCreator(mkdirAll RecursiveDirCreator, writeFile FileCreator, staticDir string) StaticFileCreator {
	return func(data *[]byte, dir, filename string) (string, error) {
		fullDir := filepath.Join(staticDir, dir)
		err := mkdirAll(fullDir, 0777)
		if err != nil {
			return "", fmt.Errorf("error while creating a new directory: %w", err)
		}
		fullPath := filepath.Join(fullDir, filename)
		err = writeFile(fullPath, *data, 0777)
		if err != nil {
			return "", fmt.Errorf("error while writing to a file: %w", err)
		}
		path := filepath.Join(dir, filename)
		return path, nil
	}
}

func NewStaticFileCreatorImpl(staticDir string) StaticFileCreator {
	return NewStaticFileCreator(os.MkdirAll, os.WriteFile, staticDir)
}
