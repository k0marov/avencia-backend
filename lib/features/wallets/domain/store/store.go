package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
)

type WalletGetter = func(userId string) (map[string]any, error)
type BalanceUpdater = func(update fs_facade.Updater, userId string, currency core.Currency, newBalance core.MoneyAmount) error
