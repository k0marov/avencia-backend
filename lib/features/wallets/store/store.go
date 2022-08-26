package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewWalletGetter(getDoc db.Getter) store.WalletGetter {
	return func(db db.DB, userId string) (map[string]any, error) {
		wallet, err := getDoc(db, "Wallets/"+userId)
		if status.Code(err) == codes.NotFound { // TODO: move such checks to the more low-level code
			return map[string]any{}, nil
		}
		if err != nil {
			return nil, core_err.Rethrow("while getting users's wallet document", err)
		}
		return wallet.Data, nil
	}
}

func NewBalanceUpdater(updDoc db.Setter) store.BalanceUpdater {
	return func(db db.DB, userId string, currency core.Currency, newBalance core.MoneyAmount) error {
		path := "Wallets/"+userId
		newValue := map[string]any{string(currency): newBalance.Num()}
		return updDoc(db, path, newValue)
	}
}
