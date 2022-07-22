package service

import "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"

// LimitChecker returns a client error if rejected; simple error if server error; nil if accepted
// LimitChecker does not update the withdrawn value, see WithdrawnUpdater
type LimitChecker = func(wantTransaction values.TransactionData) error
