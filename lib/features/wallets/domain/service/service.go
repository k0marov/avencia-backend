package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/store"
)

// BalanceGetter Should return 0 if the wallets field for the given currency is null
type BalanceGetter = func(db db.DB, userId string, currency core.Currency) (core.MoneyAmount, error)

func NewBalanceGetter(getWallet store.WalletGetter) BalanceGetter {
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
