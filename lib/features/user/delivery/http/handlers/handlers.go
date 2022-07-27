package handlers

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	"net/http"
)

func NewGetUserInfoHandler(getUserInfo service.UserInfoGetter) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		func(user auth.User, _ any) string { return user.Id },
		getUserInfo,
		func(u entities.UserInfo) api.UserInfoResponse {
			return api.UserInfoResponse{
				Id:     u.Id,
				Wallet: u.Wallet.ToResponse(),
				Limits: u.Limits.ToResponse(),
			}
		},
	)
}
