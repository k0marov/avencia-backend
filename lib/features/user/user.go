package user

import (
	"github.com/k0marov/avencia-backend/lib/features/user/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	"net/http"
)

type Handlers struct {
	GetUserInfo http.HandlerFunc
}

func NewUserHandlersImpl(getUserInfo service.UserInfoGetter) Handlers {
	return Handlers{
		GetUserInfo: handlers.NewGetUserInfoHandler(getUserInfo),
	}
}
