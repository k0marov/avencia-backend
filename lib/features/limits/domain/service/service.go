package service

import (
	"fmt"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

// LimitChecker returns a client error if rejected; simple error if server error; nil if accepted
// LimitChecker does not update the withdrawn value, see WithdrawUpdater
type LimitChecker = func(db db.DB, wantTransaction transValues.Transaction) error
type LimitsGetter = func(db db.DB, userId string) (entities.Limits, error)

type WithdrawnUpdater = func(db.DB, transValues.Transaction) error

type withdrawnUpdateGetter = func(db.DB, transValues.Transaction) (core.Money, error)


// TODO: simplify 
func NewLimitsGetter(getWithdrawns store.WithdrawsGetter, limitedCurrencies map[core.Currency]core.MoneyAmount) LimitsGetter {
	return func(db db.DB, userId string) (entities.Limits, error) {
		withdrawns, err := getWithdrawns(db, userId)
		if err != nil {
			return entities.Limits{}, core_err.Rethrow("getting current withdrawns", err)
		}
		limits := entities.Limits{}
		for curr, maxLimit := range limitedCurrencies {
			relevantWithdrawn := core.NewMoneyAmount(0)
			w := withdrawns[curr]
			if configurable.IsWithdrawLimitRelevant(w.UpdatedAt) {
				relevantWithdrawn = core.NewMoneyAmount(w.Withdrawn.Num())
			} 			
			limits[curr] = values.Limit{
				Withdrawn: relevantWithdrawn,
				Max:       maxLimit,
			}
		}
		return limits, nil
	}
}

func NewLimitChecker(getLimits LimitsGetter) LimitChecker {
	return func(db db.DB, t transValues.Transaction) error {
		if t.Money.Amount.IsPos() { // it's a deposit
			return nil
		}
		withdraw := t.Money.Amount.Neg()
		limits, err := getLimits(db, t.UserId)
		if err != nil {
			return core_err.Rethrow("while getting users's limits", err)
		}
		limit := limits[t.Money.Currency]
		if limit.Max.IsSet() && limit.Withdrawn.Add(withdraw).IsBigger(limit.Max) {
			return client_errors.WithdrawLimitExceeded
		}
		return nil
	}
}

func NewWithdrawnUpdater(getValue withdrawnUpdateGetter, update store.WithdrawUpdater) WithdrawnUpdater {
	return func(db db.DB, t transValues.Transaction) error {
		newWithdrawn, err := getValue(db, t)
		if err != nil {
			return core_err.Rethrow("getting new withdrawn value", err)
		}
		return update(db, t.UserId, newWithdrawn)
	}
}

func NewWithdrawnUpdateGetter(getLimits LimitsGetter) withdrawnUpdateGetter {
	return func(db db.DB, t transValues.Transaction) (core.Money, error) {
		if t.Money.Amount.IsPos() {
			return core.Money{}, fmt.Errorf("expected withdrawal; got deposit")
		}
		limits, err := getLimits(db, t.UserId)
		if err != nil {
			return core.Money{}, core_err.Rethrow("getting limits", err)
		}
		withdraw := t.Money.Amount.Neg()
		newWithdrawn := limits[t.Money.Currency].Withdrawn.Add(withdraw)
		return core.Money{
			Currency: t.Money.Currency,
			Amount:   newWithdrawn,
		}, nil
	}
}

