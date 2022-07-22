package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
)

// TODO: add updating limit
func NewTransactionPerformer(client firestore_facade.TransactionFirestoreFacade, updateBalance walletStore.BalanceUpdaterFactory) store.TransactionPerformer {
	return func(userId string, currency core.Currency, newBalance core.MoneyAmount) error {
		return updateBalance(client)(userId, currency, newBalance)
	}
}
