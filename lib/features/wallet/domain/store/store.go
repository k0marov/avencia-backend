package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
)

type WalletGetter = func(userId string) (map[string]any, error)
type BalanceUpdater = func(writeBatch firestore_facade.WriteBatch, userId string, currency core.Currency, newBalance core.MoneyAmount)
