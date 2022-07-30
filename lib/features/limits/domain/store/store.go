package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
)

type WithdrawsGetter = func(userId string) (fs_facade.Documents, error)

// WithdrawUpdater withdrawn's Amount should be positive
type WithdrawUpdater = func(fsUpdater fs_facade.Updater, userId string, withdrawn core.Money) error
