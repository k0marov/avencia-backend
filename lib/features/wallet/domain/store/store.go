package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
)

type WalletGetter = func(userId string) (map[string]any, error)
type BalanceUpdater = func(userId string, currency core.Currency, newBalance core.MoneyAmount) error

// BalanceUpdaterFactory is used where you need to pass a custom client (e.g inside a RunTransaction)
type BalanceUpdaterFactory = func(client firestore_facade.SimpleFirestoreFacade) BalanceUpdater
