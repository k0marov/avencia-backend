package service

import (
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

// LimitChecker returns a client error if rejected; simple error if server error; nil if accepted
// LimitChecker does not update the withdrawn value, see WithdrawUpdater
type LimitChecker = func(wantTransaction transValues.Transaction) error
type LimitsGetter = func(userId string) (entities.Limits, error)

// WithdrawnUpdateGetter computes the new Withdrawn value from a transaction;
// returns an error if transaction is not a withdrawal, in other words, when the Amount is positive
type WithdrawnUpdateGetter = func(t transValues.Transaction) (core.Money, error)

func NewLimitsGetter(getWithdrawns store.WithdrawsGetter, limitedCurrencies map[core.Currency]core.MoneyAmount) LimitsGetter {
	return func(userId string) (entities.Limits, error) {
		withdrawns, err := getWithdrawns(userId)
		if err != nil {
			return entities.Limits{}, fmt.Errorf("getting current withdrawns")
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
			return fmt.Errorf("while getting user's limits: %w", err)
		}
		limit := limits[t.Money.Currency]
		if limit.Max.IsSet() && limit.Withdrawn.Add(withdraw).IsBigger(limit.Max) {
			return client_errors.WithdrawLimitExceeded
		}
		return nil
	}
}

func NewWithdrawnUpdateGetter(getLimits LimitsGetter) WithdrawnUpdateGetter {
	return func(t transValues.Transaction) (core.Money, error) {
		if t.Money.Amount.IsPos() {
			return core.Money{}, fmt.Errorf("expected withdrawal; got deposit")
		}
		limits, err := getLimits(t.UserId)
		if err != nil {
			return core.Money{}, fmt.Errorf("getting limits: %w", err)
		}
		withdraw := t.Money.Amount.Neg()
		newWithdrawn := limits[t.Money.Currency].Withdrawn.Add(withdraw)
		return core.Money{
			Currency: t.Money.Currency,
			Amount:   newWithdrawn,
		}, nil
	}
}
