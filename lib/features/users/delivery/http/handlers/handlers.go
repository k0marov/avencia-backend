package handlers

import (
	"net/http"

	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/core/helpers/service_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/users/domain/service"
)

func NewGetUserInfoHandler(runT db.TransRunner, getUserInfo service.UserInfoGetter) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		func(user auth.User, _ *http.Request, _ http_helpers.NoJSONRequest) (string, error) { return user.Id, nil },
		service_helpers.NewDBTransService(runT, getUserInfo),
		apiResponses.UserInfoEncoder,
	)
}
