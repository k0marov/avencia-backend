package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
)

func NewWalletGetter(getDoc db.JsonGetter[entities.Wallet]) store.WalletGetter {
	return func(db db.DB, userId string) (entities.Wallet, error) {
		path := []string{"wallets", userId}
		wallet, err := getDoc(db, path)
		if core_err.IsNotFound(err) { 
			return entities.Wallet{}, nil
		}
		if err != nil {
			return nil, core_err.Rethrow("while getting users's wallet document", err)
		}
		return wallet, nil
	}
}

func NewBalanceUpdater(updDoc db.JsonUpdater[core.MoneyAmount]) store.BalanceUpdater {
	return func(db db.DB, userId string, newBalance core.Money) error {
		path := []string{"wallets", userId}
		return updDoc(db, path, string(newBalance.Currency), newBalance.Amount)
	}
}
