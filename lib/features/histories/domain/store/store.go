package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

type HistoryGetter = func(db db.TDB, userId string) ([]entities.TransEntry, error)

type TransStorer = func(db.TDB, transValues.Transaction) error 
