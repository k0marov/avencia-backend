package service

import (
	"fmt"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/store"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)


type WithdrawnUpdater = func(db.DB, transValues.Transaction) error

type withdrawnUpdateGetter = func(db.DB, transValues.Transaction) (core.Money, error)

func NewWithdrawnUpdater(getValue withdrawnUpdateGetter, update store.WithdrawUpdater) WithdrawnUpdater {
	return func(db db.DB, t transValues.Transaction) error {
		newWithdrawn, err := getValue(db, t)
		if err != nil {
			return core_err.Rethrow("getting new withdrawn value", err)
		}
		return update(db, t.UserId, newWithdrawn)
	}
}

// TODO: move usage of IsWithdrawLimitRelevant from the limits feature to the WithdrawUpdateGetter, 
// since if withdraw value is not relevant anymore, it should not only be ignored, but also reset to 0

func NewWithdrawnUpdateGetter(getWithdraws store.WithdrawsGetter) withdrawnUpdateGetter {
	return func(db db.DB, t transValues.Transaction) (core.Money, error) {
		if t.Money.Amount.IsPos() {
			return core.Money{}, fmt.Errorf("expected withdrawal; got deposit")
		}
		withdraws, err := getWithdraws(db, t.UserId)
		if err != nil {
			return core.Money{}, core_err.Rethrow("getting limits", err)
		}
		withdraw := t.Money.Amount.Neg()
		newWithdrawn := withdraws[t.Money.Currency].Withdrawn.Add(withdraw)
		return core.Money{
			Currency: t.Money.Currency,
			Amount:   newWithdrawn,
		}, nil
	}
}

