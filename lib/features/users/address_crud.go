package users

import (
	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core/crud"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
)

// TODO: remove code duplication with user_details_crud.go

func NewAddressCrudEndpoint(simpleDB db.SDB) api.Endpoint {
	store := crud.NewCRUDStore[api.Address](simpleDB, []string{"users", "addresses"})
	service := crud.Service[api.Address]{
		Store:          store,
		IgnoreNotFound: true,
		DefaultValue: api.Address{},
		IdPolicy: func(rd crud.RequestData) (id string, err error) {
			if rd.IdFromURL == "" {
				return rd.CallerId, nil
			}
			return rd.IdFromURL, nil
		},
		ReadP: crud.MustBeAuthenticated,
		WriteP: crud.MustBeAuthenticated.And(func(rd crud.RequestData) error {
			if rd.IdFromURL == rd.CallerId || rd.IdFromURL == "" {
				return nil
			}
			return client_errors.Unauthorized
		}),
	}
	handlers := crud.NewCRUDHandlers(&service)
	return crud.NewCrudEndpoint(handlers, crud.EnabledMethods{
		Read: true, Update: true,
	})
}
