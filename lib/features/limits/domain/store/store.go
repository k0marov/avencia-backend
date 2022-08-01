package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/models"
)

type WithdrawsGetter = func(userId string) ([]models.Withdrawn, error)

// WithdrawUpdater withdrawn's Amount should be positive
type WithdrawUpdater = func(fsUpdater fs_facade.Updater, userId string, withdrawn core.Money) error
