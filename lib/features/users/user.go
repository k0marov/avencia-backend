package users

import (
	"github.com/k0marov/avencia-backend/lib/features/users/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/users/domain/service"
	"net/http"
)

type Handlers struct {
	GetUserInfo http.HandlerFunc
}

func NewUserHandlersImpl(getUserInfo service.DeliveryUserInfoGetter) Handlers {
	return Handlers{
		GetUserInfo: handlers.NewGetUserInfoHandler(getUserInfo),
	}
}
