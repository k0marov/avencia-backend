package service

import (
	"fmt"
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
