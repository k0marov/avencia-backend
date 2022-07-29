package service

import (
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

// LimitChecker returns a client error if rejected; simple error if server error; nil if accepted
// LimitChecker does not update the withdrawn value, see WithdrawUpdater
type LimitChecker = func(wantTransaction transValues.Transaction) error
type LimitsGetter = func(userId string) (entities.Limits, error)

type WithdrawUpdater = func(u fs_facade.Updater, t transValues.Transaction) error

type withdrawnUpdateGetter = func(t transValues.Transaction) (core.Money, error)

func NewLimitsGetter(getWithdrawns store.WithdrawsGetter, limitedCurrencies map[core.Currency]core.MoneyAmount) LimitsGetter {
	return func(userId string) (entities.Limits, error) {
		withdrawns, err := getWithdrawns(userId)
		if err != nil {
			return entities.Limits{}, core_err.Rethrow("getting current withdrawns", err)
		}
		limits := entities.Limits{}
		for curr, maxLimit := range limitedCurrencies {
			withdrawn := withdrawns[string(curr)]
			withdrawnRelevant := core.NewMoneyAmount(0)
			if configurable.IsWithdrawLimitRelevant(withdrawn.UpdatedAt) {
				withdrawnRelevant = withdrawn.Withdrawn
			}
			limits[curr] = values.Limit{
				Withdrawn: withdrawnRelevant,
				Max:       maxLimit,
			}
		}
		return limits, nil
	}
}

func NewLimitChecker(getLimits LimitsGetter) LimitChecker {
	return func(t transValues.Transaction) error {
		if t.Money.Amount.IsPos() { // it's a deposit
			return nil
		}
		withdraw := t.Money.Amount.Neg()
		limits, err := getLimits(t.UserId)
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

func NewWithdrawnUpdateGetter(getLimits LimitsGetter) withdrawnUpdateGetter {
	return func(t transValues.Transaction) (core.Money, error) {
		if t.Money.Amount.IsPos() {
			return core.Money{}, fmt.Errorf("expected withdrawal; got deposit")
		}
		limits, err := getLimits(t.UserId)
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

func NewWithdrawUpdater(getValue withdrawnUpdateGetter, update store.WithdrawUpdater) WithdrawUpdater {
	return func(u fs_facade.Updater, t transValues.Transaction) error {
		newWithdrawn, err := getValue(t)
		if err != nil {
			return core_err.Rethrow("getting new withdrawn value", err)
		}
		err = update(u, t.UserId, newWithdrawn)
		if err != nil {
			return core_err.Rethrow("updating the withdrawn value in store", err)
		}
		return nil
	}
}
