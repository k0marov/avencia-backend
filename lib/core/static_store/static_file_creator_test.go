package static_store_test

import (
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/static_store"
)

func TestStaticFileCreator(t *testing.T) {
	staticDir := RandomString()
	tData := []byte(RandomString())
	tFile, err := core.NewFile(&tData)
	AssertNoError(t, err)
	tDir := RandomString()
	tFilename := RandomString()
	wantDir := filepath.Join(staticDir, tDir)
	wantFullPath := filepath.Join(wantDir, tFilename)
	wantPath := filepath.Join(tDir, tFilename)

	createDir := func(path string, perm fs.FileMode) error {
		return nil
	}
	t.Run("error case - creating the directory throws", func(t *testing.T) {
		recursiveDirCreator := func(path string, perm fs.FileMode) error {
			if path == wantDir && perm == 0777 {
				return RandomError()
			}
			panic("called with unexpected arguments")
		}
		_, err := static_store.NewStaticFileCreator(recursiveDirCreator, nil, staticDir)(tFile, tDir, tFilename)
		AssertSomeError(t, err)
	})
	writeFile := func(string, []byte, fs.FileMode) error {
		return nil
	}
	t.Run("error case - writing to the file throws", func(t *testing.T) {
		writeFile := func(path string, data []byte, perm fs.FileMode) error {
			if path == wantFullPath && reflect.DeepEqual(data, tData) && perm == 0777 {
				return RandomError()
			}
			panic("called with unexpected arguments")
		}
		sut := static_store.NewStaticFileCreator(createDir, writeFile, staticDir)
		_, err := sut(tFile, tDir, tFilename)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		sut := static_store.NewStaticFileCreator(createDir, writeFile, staticDir)
		gotPath, err := sut(tFile, tDir, tFilename)
		AssertNoError(t, err)
		Assert(t, gotPath, wantPath, "returned path")
	})
}
