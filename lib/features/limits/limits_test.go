package limits_test

import (
	"reflect"
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func TestLimitsGetter(t *testing.T) {
	user := RandomString()
	mockDB := NewStubDB()
	t.Run("error case - getting withdraws throws", func(t *testing.T) {
		getWithdraws := func(gotDB db.TDB, userId string) (models.Withdraws, error) {
			if gotDB == mockDB && userId == user {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		_, err := limits.NewLimitsGetter(getWithdraws, nil)(mockDB, user)
		AssertSomeError(t, err)
	})
	t.Run("forward case - forward to limits computer", func(t *testing.T) {
		withdraws := RandomWithdraws()
		tLimits := RandomLimits()
		tErr := RandomError()
		getWithdraws := func(db.TDB, string) (models.Withdraws, error) {
			return withdraws, nil
		}
		limitsComputer := func(w models.Withdraws) (limits.Limits, error) {
			if reflect.DeepEqual(w, withdraws) {
				return tLimits, tErr
			}
			panic("unexpected")
		}
		gotLimits, err := limits.NewLimitsGetter(getWithdraws, limitsComputer)(mockDB, user)
		AssertError(t, err, tErr)
		Assert(t, gotLimits, tLimits, "returned limits")
	})
}

func TestLimitsComputer(t *testing.T) {
	withdrawns := models.Withdraws{
		"BTC": {
			Withdrawn: core.NewMoneyAmount(0.001),
		},
		"RUB": {
			Withdrawn: core.NewMoneyAmount(10000),
		},
		"ETH": {
			Withdrawn: core.NewMoneyAmount(41),
		},
		"USD": {
			Withdrawn: core.NewMoneyAmount(499),
		},
	}
	limitedCurrencies := map[core.Currency]core.MoneyAmount{
		"RUB": core.NewMoneyAmount(40000),
		"USD": core.NewMoneyAmount(1000),
		"ETH": core.NewMoneyAmount(42),
		"EUR": core.NewMoneyAmount(1000),
	}
	wantLimits := limits.Limits{
		"ETH": {
			Withdrawn: core.NewMoneyAmount(41),
			Max:       core.NewMoneyAmount(42),
		},
		"USD": {
			Withdrawn: core.NewMoneyAmount(499),
			Max:       core.NewMoneyAmount(1000),
		},
		"EUR": {
			Withdrawn: core.NewMoneyAmount(0),
			Max:       core.NewMoneyAmount(1000),
		},
		"RUB": {
			Withdrawn: core.NewMoneyAmount(10000),
			Max:       core.NewMoneyAmount(40000),
		},
	}
	gotLimits, err := limits.NewLimitsComputer(limitedCurrencies)(withdrawns)
	AssertNoError(t, err)
	Assert(t, gotLimits, wantLimits, "returned limits")
}

func TestLimitChecker(t *testing.T) {
	tLimits := limits.Limits{
		"USD": limits.Limit{
			Withdrawn: core.NewMoneyAmount(500),
			Max:       core.NewMoneyAmount(600),
		},
		"RUB": limits.Limit{
			Withdrawn: core.NewMoneyAmount(400),
			Max:       core.NewMoneyAmount(10000),
		},
	}
	user := RandomString()
	mockDB := NewStubDB()

	getLimits := func(gotDB db.TDB, userId string) (limits.Limits, error) {
		if gotDB == mockDB && userId == user {
			return tLimits, nil
		}
		panic("unexpected")
	}

	t.Run("error case - getting limits throws", func(t *testing.T) {
		getLimits := func(db.TDB, string) (limits.Limits, error) {
			return nil, RandomError()
		}
		err := limits.NewLimitChecker(getLimits)(mockDB, transValues.Transaction{Money: RandomNegativeMoney()})
		AssertSomeError(t, err)
	})
	t.Run("error case - limit exceeded", func(t *testing.T) {
		trans := transValues.Transaction{
			UserId: user,
			Source: RandomTransactionSource(),
			Money: core.Money{
				Currency: "USD",
				Amount:   core.NewMoneyAmount(-200),
			},
		}
		err := limits.NewLimitChecker(getLimits)(mockDB, trans)
		AssertError(t, err, client_errors.WithdrawLimitExceeded)
	})
	t.Run("happy case", func(t *testing.T) {
		trans := transValues.Transaction{
			Source: RandomTransactionSource(),
			UserId: user,
			Money: core.Money{
				Currency: "RUB",
				Amount:   core.NewMoneyAmount(-1000),
			},
		}
		err := limits.NewLimitChecker(getLimits)(mockDB, trans)
		AssertNoError(t, err)
	})
	t.Run("happy case - value is positive (its a deposit)", func(t *testing.T) {
		trans := transValues.Transaction{
			Source: RandomTransactionSource(),
			UserId: user,
			Money: core.Money{
				Currency: "USD",
				Amount:   core.NewMoneyAmount(1000),
			},
		}
		err := limits.NewLimitChecker(getLimits)(mockDB, trans)
		AssertNoError(t, err)
	})
}
