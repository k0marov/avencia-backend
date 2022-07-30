package service

import (
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/store"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type HistoryGetter = func(userId string) ([]entities.TransEntry, error)
type TransStorer = func(fs_facade.Updater, transValues.Transaction) error


func NewHistoryGetter(getHistory store.HistoryGetter) HistoryGetter {
  return func(userId string) ([]entities.TransEntry, error) {
  	panic("unimplemented")
  }
}

func NewTransStorer() TransStorer {
  return func(u fs_facade.Updater, t transValues.Transaction) error {
    panic("unimplemented")
  }
}
