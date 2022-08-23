package handlers

import (
	"net/http"
	"net/url"

	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
)

func NewGetHistoryHandler(simpleDB db.DB, getHistory service.HistoryGetter) http.HandlerFunc {
  return http_helpers.NewAuthenticatedHandler(
		func(user auth.User, _ url.Values, _ http_helpers.NoJSONRequest) (string, error) { return user.Id, nil },
    func(userId string) ([]entities.TransEntry, error){
    	return getHistory(simpleDB, userId)
    }, 
    apiResponses.HistoryEncoder,
  ) 
}
