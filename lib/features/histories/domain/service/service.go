package service

import (
	"sort"

	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/store"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type HistoryGetter = func(userId string) ([]entities.TransEntry, error)
type TransStorer = func(fs_facade.Updater, transValues.Transaction) error


func NewHistoryGetter(getHistory store.HistoryGetter) HistoryGetter {
  return func(userId string) ([]entities.TransEntry, error) {
  	entries, err := getHistory(userId) 
  	if err != nil {
  		return []entities.TransEntry{}, core_err.Rethrow("getting history from store", err)
  	}
  	sort.Slice(entries, func(i, j int) bool {return entries[i].CreatedAt.After(entries[j].CreatedAt)})
  	return entries, nil
  }
}

func NewTransStorer(storeTrans store.TransStorer) TransStorer {
  return func(u fs_facade.Updater, t transValues.Transaction) error {
  	return storeTrans(u, t)
  }
}
