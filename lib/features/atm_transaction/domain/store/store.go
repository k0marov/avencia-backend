package store

import "github.com/k0marov/avencia-backend/lib/core"

// BalanceGetter Should return 0 if the wallet field for the given currency is null
type BalanceGetter = func(userId string, currency core.Currency) (core.MoneyAmount, error)

// BalanceUpdater also updates the withdrawn limit
type BalanceUpdater = func(userId string, currency core.Currency, newValue core.MoneyAmount) error
