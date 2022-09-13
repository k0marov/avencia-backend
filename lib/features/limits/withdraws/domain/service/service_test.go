package service_test

import (
	"testing"

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
		},
	}
	userId := RandomString()
	mockDB := NewStubDB()

	t.Run("error case - provided transaction is not a withdrawal", func(t *testing.T) {
		trans := transValues.Transaction{Money: core.Money{Amount: core.NewMoneyAmount(1000)}}
		_, err := service.NewWithdrawnUpdateGetter(nil)(mockDB, trans)
		AssertSomeError(t, err)
	})

	getWithdraws := func(db.DB, string) (models.Withdraws, error) {
		return tWithdraws, nil
	}
	t.Run("error case - getting limits throws", func(t *testing.T) {
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
	t.Run("happy case - previous withdrawn value exists", func(t *testing.T) {
		trans := transValues.Transaction{
			Source: RandomTransactionSource(),
			UserId: userId,
			Money: core.Money{
				Currency: "USD",
				Amount:   core.NewMoneyAmount(-300),
			},
		}
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(700),
		}, "returned withdrawn value")
	})
	t.Run("happy case - there is no previous withdrawn value", func(t *testing.T) {
		trans := transValues.Transaction{
			Source: RandomTransactionSource(),
			UserId: userId,
			Money: core.Money{
				Currency: "BTC",
				Amount:   core.NewMoneyAmount(-0.01),
			},
		}
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, core.Money{
			Currency: "BTC",
			Amount:   core.NewMoneyAmount(0.01),
		}, "returned withdrawn value")
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
