package service

import (
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
)

type WalletGetter = func(userId string) (entities.Wallet, error)

func NewWalletGetter(getWallet store.WalletGetter) WalletGetter {
	return func(userId string) (entities.Wallet, error) {
		storedWallet, err := getWallet(userId)
		if err != nil {
			return entities.Wallet{}, fmt.Errorf("getting wallet from store: %w", err)
		}
		wallet := map[core.Currency]core.MoneyAmount{}
		for curr, bal := range storedWallet {
			balFl, ok := bal.(float64)
			if !ok {
				return entities.Wallet{}, fmt.Errorf("balance %v for currency %v is not a float", bal, curr)
			}
			wallet[core.Currency(curr)] = core.MoneyAmount(balFl)
		}
		return wallet, nil
	}
}
