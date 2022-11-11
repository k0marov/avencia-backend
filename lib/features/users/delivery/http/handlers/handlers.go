package handlers

import (
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/users/domain/service"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	apiResponses "github.com/AvenciaLab/avencia-backend/lib/setup/api/api_responses"
)
 

func NewGetUserInfoHandler(runT db.TransRunner, getUserInfo service.UserInfoGetter) http.HandlerFunc {
  return http_helpers.NewAuthenticatedHandler(
		func(user authEntities.User, _ *http.Request, _ http_helpers.NoJSONRequest) (string, error) { return user.Id, nil },
    service_helpers.NewDBTransService(runT, getUserInfo), 
    apiResponses.UserInfoEncoder,
  ) 
}
