package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
)

func NewWalletGetter(getDoc db.JsonGetter) store.WalletGetter {
	return func(db db.DB, userId string) (map[string]any, error) {
		path := []string{"wallets", userId}
		wallet, err := getDoc(db, path)
		if core_err.IsNotFound(err) { 
			return map[string]any{}, nil
		}
		if err != nil {
			return nil, core_err.Rethrow("while getting users's wallet document", err)
		}
		return wallet.Data, nil
	}
}

func NewBalanceUpdater(updDoc db.JsonUpdater) store.BalanceUpdater {
	return func(db db.DB, userId string, currency core.Currency, newBalance core.MoneyAmount) error {
		path := []string{"wallets", userId}
		newValue := map[string]any{string(currency): newBalance.Num()}
		return updDoc(db, path, newValue)
	}
}
