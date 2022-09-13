package limits

import (
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/store"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)


type Limits map[core.Currency]Limit

type Limit struct {
	Withdrawn core.MoneyAmount
	Max       core.MoneyAmount
}

// LimitChecker returns a client error if rejected; simple error if server error; nil if accepted
// LimitChecker does not update the withdrawn value, see WithdrawUpdater
type LimitChecker = func(db db.DB, wantTransaction transValues.Transaction) error
type LimitsGetter = func(db db.DB, userId string) (Limits, error)



// TODO: simplify 
func NewLimitsGetter(getWithdraws store.WithdrawsGetter, limitedCurrencies map[core.Currency]core.MoneyAmount) LimitsGetter {
	return func(db db.DB, userId string) (Limits, error) {
		withdrawns, err := getWithdraws(db, userId)
		if err != nil {
			return Limits{}, core_err.Rethrow("getting current withdrawns", err)
		}
		limits := Limits{}
		for curr, maxLimit := range limitedCurrencies {
			w := withdrawns[curr].Withdrawn
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

