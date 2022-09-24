package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
)

type WithdrawsGetter = func(db db.TDB, userId string) (models.Withdraws, error)

// WithdrawUpdater withdrawn's Amount should be positive
type WithdrawUpdater = func(db db.TDB, userId string, withdrawn core.Money) error
