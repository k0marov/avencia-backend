package service_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
	"testing"
	"time"
)

func TestLimitsGetter(t *testing.T) {
	user := RandomString()
	t.Run("error case - getting withdrawns throws", func(t *testing.T) {
		getWithdrawns := func(userId string) (map[string]values.WithdrawnWithUpdated, error) {
			if userId == user {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewLimitsGetter(getWithdrawns, nil)(user)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		limitedCurrencies := map[core.Currency]core.MoneyAmount{
			"RUB": core.NewMoneyAmount(40000),
			"USD": core.NewMoneyAmount(1000),
			"ETH": core.NewMoneyAmount(42),
			"EUR": core.NewMoneyAmount(1000),
		}
		getWithdrawns := func(string) (map[string]values.WithdrawnWithUpdated, error) {
			return map[string]values.WithdrawnWithUpdated{
				"BTC": {Withdrawn: core.NewMoneyAmount(0.001), UpdatedAt: time.Now()},                                  // not in limited currencies
				"RUB": {Withdrawn: core.NewMoneyAmount(10000), UpdatedAt: time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC)}, // more than a year ago
				"ETH": {Withdrawn: core.NewMoneyAmount(41), UpdatedAt: time.Now().Add(-10 * time.Hour)},
				"USD": {Withdrawn: core.NewMoneyAmount(499), UpdatedAt: time.Now()},
			}, nil
		}
		limits, err := service.NewLimitsGetter(getWithdrawns, limitedCurrencies)(user)
		AssertNoError(t, err)
		wantLimits := map[core.Currency]values.Limit{
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
				Withdrawn: core.NewMoneyAmount(0),
				Max:       core.NewMoneyAmount(40000),
			},
		}
		Assert(t, limits, wantLimits, "returned limits")
	})
}

func TestLimitChecker(t *testing.T) {
	limits := entities.Limits{
		"USD": values.Limit{
			Withdrawn: core.NewMoneyAmount(500),
			Max:       core.NewMoneyAmount(600),
		},
		"RUB": values.Limit{
			Withdrawn: core.NewMoneyAmount(400),
			Max:       core.NewMoneyAmount(10000),
		},
	}
	user := RandomString()

	getLimits := func(userId string) (entities.Limits, error) {
		if userId == user {
			return limits, nil
		}
		panic("unexpected")
	}

	t.Run("error case - getting limits throws", func(t *testing.T) {
		getLimits := func(string) (entities.Limits, error) {
			return nil, RandomError()
		}
		err := service.NewLimitChecker(getLimits)(transValues.Transaction{Money: RandomNegativeMoney()})
		AssertSomeError(t, err)
	})
	t.Run("error case - limit exceeded", func(t *testing.T) {
		trans := transValues.Transaction{
			UserId: user,
			Money: core.Money{
				Currency: "USD",
				Amount:   core.NewMoneyAmount(-200),
			},
		}
		err := service.NewLimitChecker(getLimits)(trans)
		AssertError(t, err, client_errors.WithdrawLimitExceeded)
	})
	t.Run("happy case", func(t *testing.T) {
		trans := transValues.Transaction{
			UserId: user,
			Money: core.Money{
				Currency: "RUB",
				Amount:   core.NewMoneyAmount(-1000),
			},
		}
		err := service.NewLimitChecker(getLimits)(trans)
		AssertNoError(t, err)
	})
	t.Run("happy case - value is positive (its a deposit)", func(t *testing.T) {
		trans := transValues.Transaction{
			UserId: user,
			Money: core.Money{
				Currency: "USD",
				Amount:   core.NewMoneyAmount(1000),
			},
		}
		err := service.NewLimitChecker(getLimits)(trans)
		AssertNoError(t, err)
	})
}

func TestWithdrawnUpdateGetter(t *testing.T) {
	limits := map[core.Currency]values.Limit{
		"USD": {
			Withdrawn: core.NewMoneyAmount(400),
		},
	}
	userId := RandomString()

	t.Run("error case - provided transaction is not a withdrawal", func(t *testing.T) {
		_, err := service.NewWithdrawnUpdateGetter(nil)(transValues.Transaction{Money: core.Money{Amount: core.NewMoneyAmount(1000)}})
		AssertSomeError(t, err)
	})

	getLimits := func(user string) (entities.Limits, error) {
		return limits, nil
	}
	t.Run("error case - getting limits throws", func(t *testing.T) {
		getLimits := func(user string) (entities.Limits, error) {
			if user == userId {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewWithdrawnUpdateGetter(getLimits)(transValues.Transaction{UserId: userId, Money: core.Money{Amount: core.NewMoneyAmount(-1000)}})
		AssertSomeError(t, err)
	})
	t.Run("happy case - previous withdrawn value exists", func(t *testing.T) {
		trans := transValues.Transaction{
			UserId: userId,
			Money: core.Money{
				Currency: "USD",
				Amount:   core.NewMoneyAmount(-300),
			},
		}
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getLimits)(trans)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(700),
		}, "returned withdrawn value")
	})
	t.Run("happy case - there is no previous withdrawn value", func(t *testing.T) {
		trans := transValues.Transaction{
			UserId: userId,
			Money: core.Money{
				Currency: "BTC",
				Amount:   core.NewMoneyAmount(-0.01),
			},
		}
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getLimits)(trans)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, core.Money{
			Currency: "BTC",
			Amount:   core.NewMoneyAmount(0.01),
		}, "returned withdrawn value")
	})
}
