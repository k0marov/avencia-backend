package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type WalletGetter = func(db db.DB, userId string) (entities.Wallet, error)
type BalanceUpdater = func(db db.DB, userId string, newBalance core.Money) error
