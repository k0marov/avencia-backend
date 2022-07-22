package store

import "github.com/k0marov/avencia-backend/lib/core"

type TransactionPerformer = func(userId string, currency core.Currency, newValue core.MoneyAmount) error
