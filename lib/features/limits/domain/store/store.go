package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

type WithdrawnsGetter = func(userId string) (map[string]values.WithdrawnWithUpdated, error)

type WithdrawnUpdater = func(batch firestore_facade.WriteBatch, userId string, withdrawn core.Money)
