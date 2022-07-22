package wallet

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
	storeImpl "github.com/k0marov/avencia-backend/lib/features/wallet/store"
)

type Services struct {
	GetWallet             service.WalletGetter
	GetBalance            service.BalanceGetter
	BalanceUpdaterFactory store.BalanceUpdaterFactory
}

func NewWalletServicesImpl(fsClient *firestore.Client) Services {
	storeGetWallet := storeImpl.NewWalletGetter(fsClient)
	getWallet := service.NewWalletGetter(storeGetWallet)
	getBalance := service.NewBalanceGetter(getWallet)
	return Services{
		GetWallet:             getWallet,
		GetBalance:            getBalance,
		BalanceUpdaterFactory: storeImpl.NewBalanceUpdater,
	}
}
