package service

import (
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/users/domain/entities"
)

type DeliveryUserInfoGetter = func(userId string) (entities.UserInfo, error)

func NewDeliveryUserInfoGetter(simpleDB db.DB, getUserInfo UserInfoGetter) DeliveryUserInfoGetter {
	return func(userId string) (entities.UserInfo, error) {
		return getUserInfo(simpleDB, userId)
	}
}
