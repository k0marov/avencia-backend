package uploader

import (
	"errors"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
)

const FileUploadField = "file"

type UploaderFactory = func(filename string, policy Policy) http.HandlerFunc

func decodeFile(user entities.User, req *http.Request) (UserFile, error) {
	file := http_helpers.ParseFile(req, FileUploadField)
	if !file.IsSet() {
		return UserFile{}, errors.New("file could not be parsed")
	}
	return UserFile{User: user, File: file}, nil
}

// TODO: maybe test this
func NewUploaderFactory(uplService ServiceFactory) UploaderFactory {
	return func(filename string, validate Policy) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
			if !ok {
				return
			}
			uf, err := decodeFile(user, r)
			if err != nil {
				http_helpers.ThrowHTTPError(w, err)
			}
			err = uplService(validate, filename)(uf)
			if err != nil {
				http_helpers.ThrowHTTPError(w, err)
			}
		}
	}
}
