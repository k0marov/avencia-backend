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
// LimitChecker does not update the withdrawn value, see WithdrawnUpdater
type LimitChecker = func(wantTransaction transValues.Transaction) error
type LimitsGetter = func(userId string) (entities.Limits, error)

// WithdrawnUpdateGetter computes the new Withdrawn value from a transaction;
// returns an error if transaction is not a withdrawal, in other words, when the Amount is positive
type WithdrawnUpdateGetter = func(t transValues.Transaction) (core.Money, error)

func NewLimitsGetter(getWithdrawns store.WithdrawnsGetter, limitedCurrencies map[core.Currency]core.MoneyAmount) LimitsGetter {
	return func(userId string) (entities.Limits, error) {
		withdrawns, err := getWithdrawns(userId)
		if err != nil {
			return entities.Limits{}, fmt.Errorf("getting current withdrawns")
		}
		limits := entities.Limits{}
		for currStr, withdrawn := range withdrawns {
			curr := core.Currency(currStr)
			currLimit := limitedCurrencies[curr]
			if currLimit != 0 && configurable.IsWithdrawLimitRelevant(withdrawn.UpdatedAt) {
				limits[curr] = values.Limit{
					Withdrawn: withdrawn.Withdrawn,
					Max:       currLimit,
				}
			}
		}
		return limits, nil
	}
}

func NewLimitChecker(getLimits LimitsGetter) LimitChecker {
	return func(t transValues.Transaction) error {
		if t.Money.Amount > 0 { // it's a deposit
			return nil
		}
		withdraw := -t.Money.Amount
		limits, err := getLimits(t.UserId)
		if err != nil {
			return fmt.Errorf("while getting user's limits: %w", err)
		}
		limit := limits[t.Money.Currency]
		if limit.Max != 0 && limit.Withdrawn+withdraw > limit.Max {
			return client_errors.WithdrawLimitExceeded
		}
		return nil
	}
}

func NewWithdrawnUpdateGetter(getLimits LimitsGetter) WithdrawnUpdateGetter {
	return func(t transValues.Transaction) (core.Money, error) {
		if t.Money.Amount > 0 {
			return core.Money{}, fmt.Errorf("expected withdrawal; got deposit")
		}
		limits, err := getLimits(t.UserId)
		if err != nil {
			return core.Money{}, fmt.Errorf("getting limits: %w", err)
		}
		withdraw := -t.Money.Amount
		newWithdrawn := limits[t.Money.Currency].Withdrawn + withdraw
		return core.Money{
			Currency: t.Money.Currency,
			Amount:   newWithdrawn,
		}, nil
	}
}
