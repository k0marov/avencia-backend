package store

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
)

type WithdrawnUpdater = func(batch firestore_facade.WriteBatch, userId string, withdrawn core.Money)
