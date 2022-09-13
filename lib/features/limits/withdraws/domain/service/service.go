package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/config/configurable"
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

func NewWithdrawnUpdateGetter(getWithdraws store.WithdrawsGetter) withdrawnUpdateGetter {
	return func(db db.DB, t transValues.Transaction) (core.Money, error) {
		withdraws, err := getWithdraws(db, t.UserId)
		if err != nil {
			return core.Money{}, core_err.Rethrow("getting limits", err)
		}
		newWithdraw := t.Money.Amount.Neg()
		curWithdrawn := withdraws[t.Money.Currency] 
		var resWithdraw core.MoneyAmount 
		if configurable.IsWithdrawLimitRelevant(curWithdrawn.UpdatedAt) {
			resWithdraw = curWithdrawn.Withdrawn.Add(newWithdraw)
		} else {
			resWithdraw = newWithdraw
		}
		return core.Money{
			Currency: t.Money.Currency,
			Amount:   resWithdraw,
		}, nil
	}
}

