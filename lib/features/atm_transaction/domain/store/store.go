package store

// BalanceGetter Should return 0 if the wallet field for the given currency is null
type BalanceGetter = func(userId string, currency string) (float64, error)
type BalanceUpdater = func(userId, currency string, newValue float64) error
