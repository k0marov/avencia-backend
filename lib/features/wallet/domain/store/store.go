package store

type WalletGetter = func(userId string) (map[string]any, error)
