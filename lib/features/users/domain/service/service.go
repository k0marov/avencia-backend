package service

import (
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/users/domain/entities"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
)

type UserInfoGetter = func(userId string) (entities.UserInfo, error)

func NewUserInfoGetter(getWallet walletService.WalletGetter, getLimits limitsService.LimitsGetter) UserInfoGetter {
	return func(userId string) (entities.UserInfo, error) {
		wallet, err := getWallet(userId)
		if err != nil {
			return entities.UserInfo{}, core_err.Rethrow("getting wallets for users info", err)
		}
		limits, err := getLimits(userId)
		if err != nil {
			return entities.UserInfo{}, core_err.Rethrow("getting limits for users info", err)
		}
		return entities.UserInfo{
			Id:     userId,
			Wallet: wallet,
			Limits: limits,
		}, nil
	}
}