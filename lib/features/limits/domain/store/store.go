package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

type WithdrawsGetter = func(userId string) (map[string]values.WithdrawnWithUpdated, error)

// WithdrawUpdater withdrawn's Amount should be positive
type WithdrawUpdater = func(fsUpdater fs_facade.Updater, userId string, withdrawn core.Money) error
