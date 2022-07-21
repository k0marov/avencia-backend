package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
)

// SimpleFirestoreClient this interface is used instead of full firestore.Client since interfaces should be lean
type SimpleFirestoreClient interface {
	Doc(string) *firestore.DocumentRef
}

func NewBalanceGetter(client SimpleFirestoreClient) store.BalanceGetter {
	return func(userId string, currency string) (float64, error) {
		docRef := client.Doc("Wallets/" + userId)
		if docRef == nil {
			return 0, errors.New("getting document ref for user's wallet returned nil")
		}
		wallet, err := docRef.Get(context.Background())
		if err != nil {
			return 0, fmt.Errorf("while getting user's wallet: %w", err)
		}
		balance, err := wallet.DataAt(currency)
		if err != nil {
			return 0, fmt.Errorf("while getting %s field from user's wallet: %w", currency, err)
		}
		// balance is not set for this currency, default to 0
		if balance == nil {
			return 0, nil
		}
		balanceFloat, ok := balance.(float64)
		if !ok {
			return 0, fmt.Errorf("balance field (%v) is not float", balance)
		}
		return balanceFloat, nil
	}
}

func NewBalanceUpdater(client SimpleFirestoreClient) store.BalanceUpdater {
	return func(userId, currency string, newBalance float64) error {
		docRef := client.Doc("Wallets/" + userId)
		if docRef == nil {
			return errors.New("getting document ref for user's wallet returned nil")
		}
		_, err := docRef.Update(context.Background(), []firestore.Update{{Path: currency, Value: newBalance}})
		if err != nil {
			return fmt.Errorf("while updating %s with %v for %v: %w", currency, newBalance, docRef, err)
		}
		return nil
	}
}
