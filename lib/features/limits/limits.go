package limits

import (
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/store"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)


type Limits map[core.Currency]Limit

type Limit struct {
	Withdrawn core.MoneyAmount
	Max       core.MoneyAmount
}

// LimitChecker returns a client error if rejected; simple error if server error; nil if accepted
type LimitChecker = func(db db.TDB, wantTransaction transValues.Transaction) error
type LimitsGetter = func(db db.TDB, userId string) (Limits, error)



type limitsComputer = func(withdraws models.Withdraws) (Limits, error)
func NewLimitsGetter(getWithdraws store.WithdrawsGetter, compute limitsComputer) LimitsGetter {
	return func(db db.TDB, userId string) (Limits, error) {
		withdraws, err := getWithdraws(db, userId)
		if err != nil {
			return Limits{}, core_err.Rethrow("getting current withdrawns", err)
		}
		return compute(withdraws)
	}
}

func NewLimitChecker(getLimits LimitsGetter) LimitChecker {
	return func(db db.TDB, t transValues.Transaction) error {
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

func NewLimitsComputer(limitedCurrencies map[core.Currency]core.MoneyAmount) limitsComputer {
	return func(withdraws models.Withdraws) (Limits, error) {
		limits := Limits{} 
		for curr, maxLimit := range limitedCurrencies {
			w := withdraws[curr].Withdrawn
			if !w.IsSet() {
				w = core.NewMoneyAmount(0)
			}
			limits[curr] = Limit{
				Withdrawn: w,
				Max:       maxLimit,
			}
		}
		return limits, nil
	}
}
