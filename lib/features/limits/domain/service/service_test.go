package service_test

import (
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
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
			"RUB": 40000,
			"USD": 1000,
			"ETH": 42,
		}
		getWithdrawns := func(string) (map[string]values.WithdrawnWithUpdated, error) {
			return map[string]values.WithdrawnWithUpdated{
				"BTC": {Withdrawn: 0.001, UpdatedAt: time.Now()},                                  // not in limited currencies
				"RUB": {Withdrawn: 10000, UpdatedAt: time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC)}, // more than a year ago
				"ETH": {Withdrawn: 41, UpdatedAt: time.Now().Add(-10 * time.Hour)},
				"USD": {Withdrawn: 499, UpdatedAt: time.Now()},
			}, nil
		}
		limits, err := service.NewLimitsGetter(getWithdrawns, limitedCurrencies)(user)
		AssertNoError(t, err)
		wantLimits := map[core.Currency]values.Limit{
			"ETH": {
				Withdrawn: 41,
				Max:       42,
			},
			"USD": {
				Withdrawn: 499,
				Max:       1000,
			},
		}
		Assert(t, limits, wantLimits, "returned limits")
	})
}
