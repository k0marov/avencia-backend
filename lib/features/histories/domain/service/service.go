package service

import (
	"sort"

	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/store"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type DeliveryHistoryGetter = func(userId string) ([]entities.TransEntry, error) 

func NewDeliveryHistoryGetter(simpleDB db.DB, getHistory store.HistoryGetter) DeliveryHistoryGetter {
	return func(userId string) ([]entities.TransEntry, error) {
		return NewHistoryGetter(getHistory)(simpleDB, userId)
	}
}


type HistoryGetter = func(db db.DB, userId string) ([]entities.TransEntry, error)
type TransStorer = func(db.DB, transValues.Transaction) error


func NewHistoryGetter(getHistory store.HistoryGetter) HistoryGetter {
  return func(db db.DB, userId string) ([]entities.TransEntry, error) {
  	entries, err := getHistory(db, userId) 
  	if err != nil {
  		return []entities.TransEntry{}, core_err.Rethrow("getting history from store", err)
  	}
  	sort.Slice(entries, func(i, j int) bool {return entries[i].CreatedAt.After(entries[j].CreatedAt)})
  	return entries, nil
  }
}

func NewTransStorer(storeTrans store.TransStorer) TransStorer {
  return func(db db.DB, t transValues.Transaction) error {
  	return storeTrans(db, t)
  }
}
