package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
)

type TransactionPerformer = func(curBalance core.MoneyAmount, t values.Transaction) error
