package kyc

// TODO: add tests for the kyc feature

import (
	"net/http"

	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-backend/lib/core/crud"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	"github.com/go-chi/chi/v5"
)

type Status string

const (
	Unset    Status = "unset"
	Pending         = "pending"
	Verified        = "verified"
	Rejected        = "rejected"
)

type StatusModel struct {
	Status Status `json:"status"`
}

type StatusEndpointFactory = func(name string) api.Endpoint

func NewStatusEndpointFactory(simpleDB db.SDB) StatusEndpointFactory {
	return func(name string) api.Endpoint {
		store := crud.NewCRUDStore[StatusModel](simpleDB, []string{"kyc", name})
		service := crud.Service[StatusModel]{
			Store:          store,
			DefaultValue:   StatusModel{Status: Unset},
			IgnoreNotFound: true,
			IdPolicy: func(rd crud.RequestData) (id string, err error) {
				return rd.CallerId, nil
			},
			ReadP:  crud.MustBeAuthenticated,
			WriteP: crud.MustBeAuthenticated,
		}
		return func(r chi.Router) {
			r.Get("/", crud.NewCRUDHandlers(&service).Read)
			r.Patch("/", newKYCSubmitter(service))
		}
	}
}

func newKYCSubmitter(service crud.Service[StatusModel]) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		func(user entities.User, req *http.Request, _ http_helpers.NoJSONRequest) (entities.User, error) {
			return user, nil
		},
		service_helpers.NewNoResultService(func(u entities.User) error {
			return service.Update(crud.RequestData{CallerId: u.Id}, StatusModel{Status: Pending})
		}),
		http_helpers.NoResponseConverter,
	)
}
