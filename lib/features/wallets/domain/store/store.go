package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
)

type WalletGetter = func(db db.DB, userId string) (map[string]any, error)
type BalanceUpdater = func(db db.DB, userId string, currency core.Currency, newBalance core.MoneyAmount) error
