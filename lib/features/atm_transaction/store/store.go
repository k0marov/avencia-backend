package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
)

func NewBalanceGetter(getWallet service.WalletGetter) store.BalanceGetter {
	return func(userId string, currency core.Currency) (core.MoneyAmount, error) {
		wallet, err := getWallet(userId)
		if err != nil {
			return 0, fmt.Errorf("getting wallet to later extract balance: %w", err)
		}
		return wallet[currency], nil
	}
}

// TODO move balance updating logic to the wallet feature, and rename this to TransactionStorer
// TODO: add updating limit
func NewTransactionPerformer(client firestore_facade.TransactionFirestoreFacade) store.TransactionPerformer {
	return func(userId string, currency core.Currency, newBalance core.MoneyAmount) error {
		docRef := client.Doc("Wallets/" + userId)
		if docRef == nil {
			return errors.New("getting document ref for user's wallet returned nil")
		}
		_, err := docRef.Set(context.Background(), map[string]any{string(currency): float64(newBalance)}, firestore.MergeAll)
		if err != nil {
			return fmt.Errorf("while updating %s with %v for %v: %w", currency, newBalance, docRef, err)
		}
		return nil
	}
}
