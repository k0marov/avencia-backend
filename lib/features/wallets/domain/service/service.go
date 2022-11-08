package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/store"
)

type WalletCreator = func(db db.TDB, userId string, currency core.Currency) (id string, err error)
type WalletGetter = func(db db.TDB, walletId string) (entities.Wallet, error)
type WalletsGetter = func(db db.TDB, userId string) ([]entities.Wallet, error)

type BalanceUpdater = func(db db.TDB, walletId string, newBal core.MoneyAmount) error

func NewWalletCreator(create store.WalletCreator) WalletCreator {
	return func(db db.TDB, userId string, currency core.Currency) (id string, err error) {
		wallet := entities.Wallet{
			OwnerId:  userId,
			Currency: currency,
			Amount:   core.NewMoneyAmount(0),
		}
		return create(db, wallet)
	}
}

func NewWalletGetter(get store.WalletGetter) WalletGetter {
	return get
}

func NewWalletsGetter(get store.WalletsGetter) WalletsGetter {
	return get
}

func NewBalanceUpdater(update store.BalanceUpdater) BalanceUpdater {
	return update
}
