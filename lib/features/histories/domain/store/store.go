package store

import (
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
)

type HistoryGetter = func(userId string) ([]entities.TransEntry, error)
