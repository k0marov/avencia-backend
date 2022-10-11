package static_store_test

// TODO: refactor the static_store tests

import (
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/static_store"
)

func TestStaticFileCreator(t *testing.T) {
	staticDir := RandomString()
	tData := []byte(RandomString())
	tDir := RandomString()
	tFilename := RandomString()
	wantDir := filepath.Join(staticDir, tDir)
	wantFullPath := filepath.Join(wantDir, tFilename)
	wantPath := filepath.Join(tDir, tFilename)

	t.Run("should create directory", func(t *testing.T) {
		recursiveDirCreator := func(path string, perm fs.FileMode) error {
			if path == wantDir && perm == 0777 {
				return nil
			}
			panic("called with unexpected arguments")
		}
		t.Run("happy case, should write to the file", func(t *testing.T) {
			t.Run("happy case", func(t *testing.T) {
				writeFile := func(path string, data []byte, perm fs.FileMode) error {
					if path == wantFullPath && reflect.DeepEqual(data, tData) && perm == 0777 {
						return nil
					}
					panic("called with unexpected arguments")
				}
				sut := static_store.NewStaticFileCreator(recursiveDirCreator, writeFile, staticDir)
				gotPath, err := sut(&tData, tDir, tFilename)
				AssertNoError(t, err)
				Assert(t, gotPath, wantPath, "returned path")
			})
			t.Run("error case - writing to the file throws", func(t *testing.T) {
				writeFile := func(string, []byte, fs.FileMode) error {
					return RandomError()
				}
				sut := static_store.NewStaticFileCreator(recursiveDirCreator, writeFile, staticDir)
				_, err := sut(&tData, tDir, tFilename)
				AssertSomeError(t, err)
			})
		})
		t.Run("error case - mkdirAll throws", func(t *testing.T) {
			recursiveDirCreator := func(path string, perm fs.FileMode) error {
				return RandomError()
			}
			sut := static_store.NewStaticFileCreator(recursiveDirCreator, nil, staticDir) 
			_, err := sut(&tData, tDir, tFilename)
			AssertSomeError(t, err)
		})
	})
}
