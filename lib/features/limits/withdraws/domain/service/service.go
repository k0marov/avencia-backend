package service

import (
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/store"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	"github.com/AvenciaLab/avencia-backend/lib/setup/config/configurable"
)


type WithdrawnUpdater = func(db.TDB, transValues.Transaction) error

type withdrawnUpdateGetter = func(db.TDB, transValues.Transaction) (core.Money, error)

func NewWithdrawnUpdater(getValue withdrawnUpdateGetter, update store.WithdrawUpdater) WithdrawnUpdater {
	return func(db db.TDB, t transValues.Transaction) error { 
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
	return func(db db.TDB, t transValues.Transaction) (core.Money, error) {
		withdraws, err := getWithdraws(db, t.UserId)
		if err != nil {
			return core.Money{}, core_err.Rethrow("getting limits", err)
		}
		curWithdrawn := withdraws[t.Money.Currency] 
		newWithdraw := t.Money.Amount.Neg()
    
    result := core.Money{Currency: t.Money.Currency}
		if configurable.IsWithdrawLimitRelevant(time.Unix(curWithdrawn.UpdatedAt, 0)) {
			result.Amount = curWithdrawn.Withdrawn.Add(newWithdraw)
		} else {
			result.Amount = newWithdraw
		}
		return result, nil
	}
}

