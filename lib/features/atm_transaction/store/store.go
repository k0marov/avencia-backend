package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
)

func NewTransactionPerformer(client *firestore.Client, updateBalance walletStore.BalanceUpdater, getNewWithdrawn limitsService.WithdrawnUpdateGetter, updateWithdrawn limitsStore.WithdrawnUpdater) store.TransactionPerformer {
	return func(curBal core.MoneyAmount, t values.Transaction) error {
		batch := client.Batch()
		if t.Money.Amount < 0 {
			withdrawn, err := getNewWithdrawn(t)
			if err != nil {
				return fmt.Errorf("getting the new 'withdrawn' value: %w", err)
			}
			updateWithdrawn(batch, t.UserId, withdrawn)
		}
		updateBalance(batch, t.UserId, t.Money.Currency, curBal+t.Money.Amount)
		_, err := batch.Commit(context.Background())
		return err
	}
}
