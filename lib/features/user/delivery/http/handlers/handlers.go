package handlers

import (
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	"net/http"
)

func NewGetUserInfoHandler(getUserInfo service.UserInfoGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		userInfo, err := getUserInfo(user.Id)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, userInfo.ToResponse())
	}
}
