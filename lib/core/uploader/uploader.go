package uploader

import (
	"errors"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/static_store"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
)

const FileUploadField = "file"

type UploaderFactory = func(filename string) http.HandlerFunc

type UserFile struct {
	User entities.User
	File core.File
}

func fileDecoder(user entities.User, req *http.Request, _ http_helpers.NoJSONRequest) (UserFile, error) {
	file := http_helpers.ParseFile(req, FileUploadField)
	if !file.IsSet() {
		return UserFile{}, errors.New("file could not be parsed")
	}
	return UserFile{User: user, File: file}, nil
}

func NewUploaderFactory(createFile static_store.StaticFileCreator) UploaderFactory {
	return func(filename string) http.HandlerFunc {
		return http_helpers.NewAuthenticatedHandler(
			fileDecoder,
			service_helpers.NewNoResultService(func(uf UserFile) error {
				_, err := createFile(uf.File, uf.User.Id, filename) 
				return err
			}),
			http_helpers.NoResponseConverter,
		)
	}
}
