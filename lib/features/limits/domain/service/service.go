package service

import (
	"github.com/k0marov/avencia-backend/lib/core"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
)

// LimitChecker returns a client error if rejected; simple error if server error; nil if accepted
// LimitChecker does not update the withdrawn value, see WithdrawnUpdater
type LimitChecker = func(wantTransaction transValues.TransactionData) error
type LimitsGetter = func(userId string) (entities.Limits, error)

// WithdrawnUpdateGetter computes the new Withdrawn value from a transaction;
// returns an error if transaction is not a withdrawal, in other words, when the Amount is positive
type WithdrawnUpdateGetter = func(t transValues.TransactionData) (core.Money, error)
