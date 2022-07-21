package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewWalletGetter(client firestore_facade.SimpleFirestoreFacade) store.WalletGetter {
	return func(userId string) (map[string]any, error) {
		docRef := client.Doc("Wallets/" + userId)
		if docRef == nil {
			return nil, errors.New("getting document ref for user's wallet returned nil")
		}
		wallet, err := docRef.Get(context.Background())
		if status.Code(err) == codes.NotFound {
			return map[string]any{}, nil
		}
		if err != nil {
			return nil, fmt.Errorf("while getting user's wallet document: %w", err)
		}
		return wallet.Data(), nil
	}
}
