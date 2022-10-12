package uploader

import (
	"errors"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
)

const FileUploadField = "file"

type UploaderFactory = func(filename string, policy Policy) http.HandlerFunc

func fileDecoder(user entities.User, req *http.Request, _ http_helpers.NoJSONRequest) (UserFile, error) {
	file := http_helpers.ParseFile(req, FileUploadField)
	if !file.IsSet() {
		return UserFile{}, errors.New("file could not be parsed")
	}
	return UserFile{User: user, File: file}, nil
}

func NewUploaderFactory(service ServiceFactory) UploaderFactory {
	return func(filename string, validate Policy) http.HandlerFunc {
		return http_helpers.NewAuthenticatedHandler(
			fileDecoder,
			service_helpers.NewNoResultService(service(validate, filename)),
			http_helpers.NoResponseConverter,
		)
	}
}
