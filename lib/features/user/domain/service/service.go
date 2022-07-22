package service

import (
	"fmt"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
)

type UserInfoGetter = func(userId string) (entities.UserInfo, error)

// TODO: add returning the remaining limits

func NewUserInfoGetter(getWallet walletService.WalletGetter) UserInfoGetter {
	return func(userId string) (entities.UserInfo, error) {
		wallet, err := getWallet(userId)
		if err != nil {
			return entities.UserInfo{}, fmt.Errorf("getting wallet for user info: %w", err)
		}
		return entities.UserInfo{Id: userId, Wallet: wallet}, nil
	}
}
