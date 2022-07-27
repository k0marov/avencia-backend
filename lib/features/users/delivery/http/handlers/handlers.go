package handlers

import (
	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/users/domain/service"
	"net/http"
	"net/url"
)

func NewGetUserInfoHandler(getUserInfo service.UserInfoGetter) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		func(user auth.User, _ url.Values, _ http_helpers.NoJSONRequest) (string, error) { return user.Id, nil },
		getUserInfo,
		apiResponses.UserInfoEncoder,
	)
}
