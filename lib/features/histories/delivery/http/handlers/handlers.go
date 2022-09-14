package handlers

import (
	"net/http"

	apiResponses "github.com/AvenciaLab/avencia-backend/lib/setup/api/api_responses"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/service"
)

func NewGetHistoryHandler(runT db.TransRunner, getHistory service.HistoryGetter) http.HandlerFunc {
  return http_helpers.NewAuthenticatedHandler(
		func(user authEntities.User, _ *http.Request, _ http_helpers.NoJSONRequest) (string, error) { return user.Id, nil },
		service_helpers.NewDBTransService(runT, getHistory),
    apiResponses.HistoryEncoder,
  ) 
}
