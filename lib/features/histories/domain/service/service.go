package service

import (
	transValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
)

type HistoryGetter = func(userId string) ([]entities.TransEntry, error)
type TransStorer = func(transValues.Transaction) error


func NewHistoryGetter() HistoryGetter {
  return func(userId string) ([]entities.TransEntry, error) {
    panic("unimplemented")
  }
}

func NewTransStorer() TransStorer {
  return func(t transValues.Transaction) error {
    panic("unimplemented")
  }
}
