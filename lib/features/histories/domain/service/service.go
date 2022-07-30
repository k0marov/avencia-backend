package service

import (
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
)

type HistoryGetter = func(userId string) ([]entities.TransEntry, error)
type TransStorer = func(fs_facade.Updater, transValues.Transaction) error


func NewHistoryGetter() HistoryGetter {
  return func(userId string) ([]entities.TransEntry, error) {
    panic("unimplemented")
  }
}

func NewTransStorer() TransStorer {
  return func(u fs_facade.Updater, t transValues.Transaction) error {
    panic("unimplemented")
  }
}
