package service_test

import (
	"testing"
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/service"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func TestWithdrawnUpdateGetter(t *testing.T) {
	tWithdraws := models.Withdraws{
		"USD": {
			Withdrawn: core.NewMoneyAmount(400),
			UpdatedAt: time.Now(), // still relevant
		},
		"RUB": {
			Withdrawn: core.NewMoneyAmount(1000),
			UpdatedAt: TimeWithYear(2000), // no longer relevant
		},
	}
	userId := RandomString()
	mockDB := NewStubDB()

	getWithdraws := func(db.DB, string) (models.Withdraws, error) {
		return tWithdraws, nil
	}

	t.Run("error case - getting withdraws throws", func(t *testing.T) {
		getWithdraws := func(gotDB db.DB, user string) (models.Withdraws, error) {
			if gotDB == mockDB && user == userId {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		trans := transValues.Transaction{UserId: userId, Money: core.Money{Amount: core.NewMoneyAmount(-1000)}}
		_, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, trans)
		AssertSomeError(t, err)
	})
	getTestTrans := func(w core.Money) transValues.Transaction {
		return transValues.Transaction{
			Source: RandomTransactionSource(),
			UserId: userId,
			Money:  w,
		}
	}
	negate := func(w core.Money) core.Money {
		return core.Money{
			Currency: w.Currency,
			Amount:   w.Amount.Neg(),
		}
	}
	t.Run("happy case - previous withdrawn value exists", func(t *testing.T) {
		trans := getTestTrans(core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(-300),
		})
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(700),
		}, "returned withdrawn value")
	})
	t.Run("happy case - there is no previous withdrawn value", func(t *testing.T) {
		w := core.Money{
			Currency: "BTC",
			Amount:   core.NewMoneyAmount(-0.01),
		}
		trans := getTestTrans(w)
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, negate(w), "returned withdrawn value")
	})
	t.Run("happy case - previous withdrawn value is not relevant anymore (according to configurable.IsWithdrawLimitRelevant)", func(t *testing.T) {
		w := core.Money{
			Currency: "RUB",
			Amount:   core.NewMoneyAmount(420),
		}
		trans := getTestTrans(w)
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, negate(w), "returned withdraw value") 
	})
}

func TestWithdrawnUpdater(t *testing.T) {
	trans := transValues.Transaction{
		Source: RandomTransactionSource(),
		UserId: RandomString(),
		Money:  RandomNegativeMoney(),
	}
	newWithdrawn := RandomPositiveMoney()
	mockDB := NewStubDB()
	t.Run("error case - getting new value throws", func(t *testing.T) {
		getValue := func(gotDB db.DB, gotTrans transValues.Transaction) (core.Money, error) {
			if gotDB == mockDB && gotTrans == trans {
				return core.Money{}, RandomError()
			}
			panic("unexpected")
		}
		err := service.NewWithdrawnUpdater(getValue, nil)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("forward case", func(t *testing.T) {
		tErr := RandomError()
		getValue := func(db.DB, transValues.Transaction) (core.Money, error) {
			return newWithdrawn, nil
		}
		update := func(gotDB db.DB, userId string, value core.Money) error {
			if gotDB == mockDB && userId == trans.UserId && value == newWithdrawn {
				return tErr
			}
			panic("unexpected")
		}
		err := service.NewWithdrawnUpdater(getValue, update)(mockDB, trans)
		AssertError(t, err, tErr)

	})
}
