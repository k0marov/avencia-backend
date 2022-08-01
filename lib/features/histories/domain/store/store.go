package store

import (
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type HistoryGetter = func(userId string) ([]entities.TransEntry, error)

type TransStorer = func(fs_facade.Updater, transValues.Transaction) error 
