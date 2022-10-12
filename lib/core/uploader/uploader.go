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

type UploaderFactory = func(field string, filename string) http.HandlerFunc

type UserFile struct {
	User entities.User
	File core.File
}

func NewUploaderFactory(createFile static_store.StaticFileCreator) UploaderFactory {
	return func(field string, filename string) http.HandlerFunc {
		return http_helpers.NewAuthenticatedHandler(
			func(user entities.User, req *http.Request, _ http_helpers.NoJSONRequest) (UserFile, error) {
				file := http_helpers.ParseFile(req, field)
				if !file.IsSet() {
					return UserFile{}, errors.New("file could not be parsed")
				}
				return UserFile{User: user, File: file}, nil
			},
			func(uf UserFile) (service_helpers.Nothing, error) {
				_, err := createFile(uf.File, uf.User.Id, filename)
				return service_helpers.Nothing{}, err
			},
			http_helpers.NoResponseConverter,
		)
	}
}
