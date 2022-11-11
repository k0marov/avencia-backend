package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	authStore "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/store"
	limitsService "github.com/AvenciaLab/avencia-backend/lib/features/limits"
	"github.com/AvenciaLab/avencia-backend/lib/features/users/domain/entities"
	wallets "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
)


type UserInfoGetter = func(db db.TDB, userId string) (entities.UserInfo, error)


// TODO: add getting History to the UserInfoGetter

func NewUserInfoGetter(
	getWallets wallets.WalletsGetter, 
	getLimits limitsService.LimitsGetter, 
	getUser authStore.UserGetter,
) UserInfoGetter {
	return func(db db.TDB, userId string) (entities.UserInfo, error) {
		wallets, err := getWallets(db, userId)
		if err != nil {
			return entities.UserInfo{}, core_err.Rethrow("getting wallets for users info", err)
		}
		limits, err := getLimits(db, userId)
		if err != nil {
			return entities.UserInfo{}, core_err.Rethrow("getting limits for users info", err)
		}
		user, err := getUser(userId) 
		if err != nil {
			return entities.UserInfo{}, core_err.Rethrow("getting detailed user info", err)
		}
		return entities.UserInfo{
			User: user,
			Wallets: wallets,
			Limits: limits,
		}, nil
	}
}
