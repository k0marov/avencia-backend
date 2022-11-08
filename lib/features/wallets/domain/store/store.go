package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type BalanceUpdater = func(db db.TDB, walletId string, newBalance core.MoneyAmount) error

type WalletCreator = func(db.TDB, entities.Wallet) (id string, err error)
type WalletGetter = func(db db.TDB, walletId string) (entities.Wallet, error) 
type WalletsGetter = func(db db.TDB, userId string) ([]entities.Wallet, error)

