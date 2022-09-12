package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

type HistoryGetter = func(db db.DB, userId string) ([]entities.TransEntry, error)

type TransStorer = func(db.DB, transValues.Transaction) error 
