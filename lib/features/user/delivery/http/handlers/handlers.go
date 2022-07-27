package handlers

import (
	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	"net/http"
)

func NewGetUserInfoHandler(getUserInfo service.UserInfoGetter) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		func(user auth.User, _ any) string { return user.Id },
		getUserInfo,
		apiResponses.UserInfoEncoder,
	)
}
