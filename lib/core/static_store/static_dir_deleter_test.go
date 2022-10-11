package static_store_test

import (
	"path/filepath"
	"testing"

	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/static_store"
)

func TestStaticDirDeleter(t *testing.T) {
	tPath := RandomString()
	staticDir := RandomString()
	wantDirPath := filepath.Join(staticDir, tPath)

	t.Run("error case - deleting the dir throws", func(t *testing.T) {
		deleteDir := func(dir string) error {
			if dir == wantDirPath {
				return RandomError()
			}
			panic("unexpected")
		}
		err := static_store.NewStaticDirDeleter(deleteDir, staticDir)(tPath)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		deleteDir := func(dir string) error {
			return nil
		}
		err := static_store.NewStaticDirDeleter(deleteDir, staticDir)(tPath)
		AssertNoError(t, err)
	})
}
