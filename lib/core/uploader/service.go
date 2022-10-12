package uploader


// TODO: add tests for the uploader feature 

import (
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/static_store"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
)

type UserFile struct {
	User entities.User
	File core.File
}

const MaxFileSize = 6 * 1000 * 1000 

type Policy = func(UserFile) (_ client_errors.ClientError, ok bool)

func SimpleSizePolicy(uf UserFile) (_ client_errors.ClientError, ok bool) {
  d, err := uf.File.Data()
  if err != nil {
    return client_errors.InvalidFile, false
  }
  if len(d) > MaxFileSize { 
		return client_errors.FileTooBig, false
  }
  return client_errors.ClientError{}, true 
}

type Service = func(UserFile) error
type ServiceFactory = func(p Policy, filename string) Service

func NewServiceFactory(createFile static_store.StaticFileCreator) ServiceFactory {
	return func(up Policy, filename string) Service {
		return func(uf UserFile) error {
		  if clErr, ok := up(uf); ok != true {
        return clErr
		  }
			_, err := createFile(uf.File, uf.User.Id, filename)
			return err
		}
	}
}
