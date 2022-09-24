package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	authStore "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/store"
	limitsService "github.com/AvenciaLab/avencia-backend/lib/features/limits"
	"github.com/AvenciaLab/avencia-backend/lib/features/users/domain/entities"
	walletStore "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/store"
)


type UserInfoGetter = func(db db.TDB, userId string) (entities.UserInfo, error)

func NewUserInfoGetter(
	getWallet walletStore.WalletGetter, 
	getLimits limitsService.LimitsGetter, 
	getUser authStore.UserGetter,
) UserInfoGetter {
	return func(db db.TDB, userId string) (entities.UserInfo, error) {
		wallet, err := getWallet(db, userId)
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
			Wallet: wallet,
			Limits: limits,
		}, nil
	}
}
