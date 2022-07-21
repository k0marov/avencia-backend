package wallet

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/wallet/store"
)

func NewWalletGetterImpl(fsClient *firestore.Client) service.WalletGetter {
	return service.NewWalletGetter(store.NewWalletGetter(fsClient))
}
