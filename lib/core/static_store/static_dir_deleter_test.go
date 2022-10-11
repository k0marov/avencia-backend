package static_store_test

import (
	"path/filepath"
	"testing"

	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/static_store"
)

func TestStaticDirDeleter(t *testing.T) {
	tPath := RandomString()
	wantDirPath := filepath.Join(static_store.StaticDir, tPath)
	t.Run("happy case", func(t *testing.T) {
		deleteDir := func(dir string) error {
			if dir == wantDirPath {
				return nil
			}
			panic("unexpected args")
		}
		sut := static_store.NewStaticDirDeleter(deleteDir)
		err := sut(tPath)
		AssertNoError(t, err)
	})
	t.Run("error case - deleting the dir throws", func(t *testing.T) {
		deleteDir := func(string) error {
			return RandomError()
		}
		err := static_store.NewStaticDirDeleter(deleteDir)(tPath)
		AssertSomeError(t, err)
	})
}
