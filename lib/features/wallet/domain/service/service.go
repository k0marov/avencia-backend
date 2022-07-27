package service

import (
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
)

type WalletGetter = func(userId string) (entities.Wallet, error)

// BalanceGetter Should return 0 if the wallet field for the given currency is null
type BalanceGetter = func(userId string, currency core.Currency) (core.MoneyAmount, error)

func NewWalletGetter(getWallet store.WalletGetter) WalletGetter {
	return func(userId string) (entities.Wallet, error) {
		storedWallet, err := getWallet(userId)
		if err != nil {
			return entities.Wallet{}, core_err.Rethrow("getting wallet from store", err)
		}
		wallet := map[core.Currency]core.MoneyAmount{}
		for curr, bal := range storedWallet {
			balFl, ok := bal.(float64)
			if !ok {
				return entities.Wallet{}, fmt.Errorf("balance %v for currency %v is not a float", bal, curr)
			}
			wallet[core.Currency(curr)] = core.NewMoneyAmount(balFl)
		}
		return wallet, nil
	}
}

func NewBalanceGetter(getWallet WalletGetter) BalanceGetter {
	return func(userId string, currency core.Currency) (core.MoneyAmount, error) {
		wallet, err := getWallet(userId)
		if err != nil {
			return core.NewMoneyAmount(0), core_err.Rethrow("getting wallet to later extract balance", err)
		}
		bal := wallet[currency]
		if !bal.IsSet() {
			return core.NewMoneyAmount(0), nil
		}
		return bal, nil
	}
}
