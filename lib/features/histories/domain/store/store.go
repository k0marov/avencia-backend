package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
)

type HistoryGetter = func(db db.TDB, userId string) ([]entities.TransEntry, error)

type EntryStorer = func(db db.TDB, userId string, entry entities.TransEntry) error 
