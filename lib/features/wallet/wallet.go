package wallet

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
	storeImpl "github.com/k0marov/avencia-backend/lib/features/wallet/store"
)

type Services struct {
	GetWallet      service.WalletGetter
	GetBalance     service.BalanceGetter
	BalanceUpdater store.BalanceUpdater
}

// TODO: write some integration tests (later)

func NewWalletServicesImpl(fsClient *firestore.Client) Services {
	walletDocGetter := storeImpl.NewWalletDocGetter(firestore_facade.NewDocGetter(fsClient))

	storeGetWallet := storeImpl.NewWalletGetter(walletDocGetter)
	updateBalance := storeImpl.NewBalanceUpdater(walletDocGetter)

	getWallet := service.NewWalletGetter(storeGetWallet)
	getBalance := service.NewBalanceGetter(getWallet)
	return Services{
		GetWallet:      getWallet,
		GetBalance:     getBalance,
		BalanceUpdater: updateBalance,
	}
}
