package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WalletDocGetter the passed in userId shouldn't be empty
type WalletDocGetter = func(userId string) *firestore.DocumentRef

func NewWalletDocGetter(client firestore_facade.Simple) WalletDocGetter {
	return func(userId string) *firestore.DocumentRef {
		doc := client.Doc("Wallets/" + userId)
		if doc == nil {
			panic("getting document ref for user's wallet returned nil. Probably userId is empty.")
		}
		return doc
	}
}

func NewWalletGetter(getWalletDoc WalletDocGetter) store.WalletGetter {
	return func(userId string) (map[string]any, error) {
		wallet, err := getWalletDoc(userId).Get(context.Background())
		if status.Code(err) == codes.NotFound {
			return map[string]any{}, nil
		}
		if err != nil {
			return nil, fmt.Errorf("while getting user's wallet document: %w", err)
		}
		return wallet.Data(), nil
	}
}

func NewBalanceUpdater(getWalletDoc WalletDocGetter) store.BalanceUpdater {
	return func(batch firestore_facade.WriteBatch, userId string, currency core.Currency, newBalance core.MoneyAmount) {
		doc := getWalletDoc(userId)
		newValue := map[string]any{string(currency): newBalance.Num()}
		batch.Set(doc, newValue, firestore.MergeAll)
	}
}
