package service

import (
	"sort"

	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/store"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

type HistoryGetter = func(db db.TDB, userId string) ([]entities.TransEntry, error)
type TransStorer = func(db.TDB, transValues.Transaction) error


func NewHistoryGetter(getHistory store.HistoryGetter) HistoryGetter {
  return func(db db.TDB, userId string) ([]entities.TransEntry, error) {
  	e, err := getHistory(db, userId) 
  	if err != nil {
  		return []entities.TransEntry{}, core_err.Rethrow("getting history from store", err)
  	}
  	sort.Slice(e, func(i, j int) bool {return e[i].CreatedAt > (e[j].CreatedAt)})
  	return e, nil
  }
}

func NewTransStorer(storeTrans store.TransStorer) TransStorer {
  return func(db db.TDB, t transValues.Transaction) error {
  	return storeTrans(db, t)
  }
}
