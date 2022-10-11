package static_store

import (
	"fmt"
	"os"
	"path/filepath"
)

// DirDeleter os.RemoveAll implements this
type DirDeleter = func(dir string) error

func NewStaticDirDeleter(deleteDir DirDeleter) StaticDirDeleter {
	return func(dir string) error {
		fullDir := filepath.Join(StaticDir, dir)
		err := deleteDir(fullDir)
		if err != nil {
			return fmt.Errorf("while deleting a static dir (%v) : %w", fullDir, err)
		}
		return nil
	}
}

func NewStaticDirDeleterImpl() StaticDirDeleter {
	return NewStaticDirDeleter(os.RemoveAll)
}
