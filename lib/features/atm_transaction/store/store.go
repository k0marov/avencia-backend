package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
)

// TODO: add updating limit
func NewTransactionPerformer(client *firestore.Client, updateBalance walletStore.BalanceUpdater) store.TransactionPerformer {
	return func(curBal core.MoneyAmount, t values.TransactionData) error {
		batch := client.Batch()
		updateBalance(batch, t.UserId, t.Money.Currency, curBal+t.Money.Amount)
		_, err := batch.Commit(context.Background())
		return err
	}
}
