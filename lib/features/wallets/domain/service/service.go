package service

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
)

type WalletGetter = func(db db.DB, userId string) (entities.Wallet, error)

// BalanceGetter Should return 0 if the wallets field for the given currency is null
type BalanceGetter = func(db db.DB, userId string, currency core.Currency) (core.MoneyAmount, error)

func NewWalletGetter(getWallet store.WalletGetter) WalletGetter {
	return func(db db.DB, userId string) (entities.Wallet, error) {
		storedWallet, err := getWallet(db, userId)
		if err != nil {
			return entities.Wallet{}, core_err.Rethrow("getting wallets from store", err)
		}
		wallet := map[core.Currency]core.MoneyAmount{}
		for curr, bal := range storedWallet {
			balFl, err := general_helpers.DecodeFloat(bal)
			if err != nil {
				return entities.Wallet{}, core_err.Rethrow("decoding balance", err)
			}
			wallet[core.Currency(curr)] = core.NewMoneyAmount(balFl)
		}
		return wallet, nil
	}
}

func NewBalanceGetter(getWallet WalletGetter) BalanceGetter {
	return func(db db.DB, userId string, currency core.Currency) (core.MoneyAmount, error) {
		wallet, err := getWallet(db, userId)
		if err != nil {
			return core.NewMoneyAmount(0), core_err.Rethrow("getting wallets to later extract balance", err)
		}
		bal := wallet[currency]
		if !bal.IsSet() {
			return core.NewMoneyAmount(0), nil
		}
		return bal, nil
	}
}
