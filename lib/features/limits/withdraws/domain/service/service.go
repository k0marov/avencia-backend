package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/setup/config/configurable"
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
		if t.Money.Amount.IsPos() { // it is a deposit - no update needed
			return nil
		}
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
		curWithdrawn := withdraws[t.Money.Currency] 
		newWithdraw := t.Money.Amount.Neg()
    
    result := core.Money{Currency: t.Money.Currency}
		if configurable.IsWithdrawLimitRelevant(curWithdrawn.UpdatedAt) {
			result.Amount = curWithdrawn.Withdrawn.Add(newWithdraw)
		} else {
			result.Amount = newWithdraw
		}
		return result, nil
	}
}

