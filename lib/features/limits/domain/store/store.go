package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/domain/models"
)

type WithdrawsGetter = func(db db.DB, userId string) (models.Withdraws, error)

// WithdrawUpdater withdrawn's Amount should be positive
type WithdrawUpdater = func(db db.DB, userId string, withdrawn core.Money) error
