package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
)

// TODO: rename TransactionData to Transaction

type TransactionPerformer = func(curBalance core.MoneyAmount, t values.TransactionData) error
