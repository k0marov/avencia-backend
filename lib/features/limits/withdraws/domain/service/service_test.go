package service_test

import (
	"testing"
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/service"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestWithdrawnUpdateGetter(t *testing.T) {
	tWithdraws := models.Withdraws{
		"USD": {
			Withdrawn: core.NewMoneyAmount(400),
			UpdatedAt: time.Now().Unix(), // still relevant
		},
		"RUB": {
			Withdrawn: core.NewMoneyAmount(1000),
			UpdatedAt: TimeWithYear(2000).Unix(), // no longer relevant
		},
	}
	userId := RandomString()
	mockDB := NewStubDB()

	getWithdraws := func(db.TDB, string) (models.Withdraws, error) {
		return tWithdraws, nil
	}

	t.Run("error case - getting withdraws throws", func(t *testing.T) {
		getWithdraws := func(gotDB db.TDB, user string) (models.Withdraws, error) {
			if gotDB == mockDB && user == userId {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		money := core.Money{Amount: core.NewMoneyAmount(-1000)}
		_, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, userId, money)
		AssertSomeError(t, err)
	})
	t.Run("happy case - previous withdrawn value exists", func(t *testing.T) {
		money := core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(-300),
		}
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, userId, money)
		AssertNoError(t, err)
		Assert(t, newWithdrawn, core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(700),
		}, "returned withdrawn value")
	})
	t.Run("happy case - there is no previous withdrawn value", func(t *testing.T) {
		money := core.Money{
			Currency: "BTC",
			Amount:   core.NewMoneyAmount(-0.01),
		}
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, userId, money)
		AssertNoError(t, err)
		expected := core.Money{
			Currency: "BTC",
			Amount:   core.NewMoneyAmount(0.01),
		}
		Assert(t, newWithdrawn, expected, "returned withdrawn value")
	})
	t.Run("happy case - previous withdrawn value is not relevant anymore (according to configurable.IsWithdrawLimitRelevant)", func(t *testing.T) {
		money := core.Money{
			Currency: "RUB",
			Amount:   core.NewMoneyAmount(-420),
		}
		newWithdrawn, err := service.NewWithdrawnUpdateGetter(getWithdraws)(mockDB, userId, money)
		AssertNoError(t, err)
		expected := core.Money{
			Currency: "RUB",
			Amount:   core.NewMoneyAmount(420),
		}
		Assert(t, newWithdrawn, expected, "returned withdraw value")
	})
}

func TestWithdrawnUpdater(t *testing.T) {
	userId := RandomString()
	money := RandomNegativeMoney()
	newWithdrawn := RandomPositiveMoney()

	mockDB := NewStubDB()
	t.Run("early return case - the provided transaction is not a withdrawal", func(t *testing.T) {
		err := service.NewWithdrawnUpdater(nil, nil)(mockDB, userId, RandomPositiveMoney())
		AssertNoError(t, err)
	})
	t.Run("error case - getting new value throws", func(t *testing.T) {
		getValue := func(gotDB db.TDB, user string, gotMoney core.Money) (core.Money, error) {
			if gotDB == mockDB && user == userId && gotMoney == money {
				return core.Money{}, RandomError()
			}
			panic("unexpected")
		}
		err := service.NewWithdrawnUpdater(getValue, nil)(mockDB, userId, money)
		AssertSomeError(t, err)
	})
	t.Run("forward case", func(t *testing.T) {
		tErr := RandomError()
		getValue := func(db.TDB, string, core.Money) (core.Money, error) {
			return newWithdrawn, nil
		}
		update := func(gotDB db.TDB, userId string, value core.Money) error {
			if gotDB == mockDB && userId == userId && value == newWithdrawn {
				return tErr
			}
			panic("unexpected")
		}
		err := service.NewWithdrawnUpdater(getValue, update)(mockDB, userId, money)
		AssertError(t, err, tErr)

	})
}

func TestTransWithdrawnUpdater(t *testing.T) {
	mockDB := NewStubDB()
	trans := RandomTransactionData()
	wallet := RandomWallet()
	getWallet := func(gotDB db.TDB, walletId string) (wEntities.Wallet, error) {
		if gotDB == mockDB && walletId == trans.WalletId {
			return wallet, nil
		}
		panic("unexpected")
	}
	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(gotDB db.TDB, walletId string) (wEntities.Wallet, error) {
			return wEntities.Wallet{}, RandomError()
		}
		err := service.NewTransWithdrawnUpdater(getWallet, nil)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("forward case - forward to updateWithdrawn", func(t *testing.T) {
		tErr := RandomError()
		upd := func(gotDB db.TDB, userId string, money core.Money) error {
			if gotDB == mockDB && userId == wallet.OwnerId &&
				money == trans.Money {
				return tErr
			}
			panic("unexpected")
		}
		err := service.NewTransWithdrawnUpdater(getWallet, upd)(mockDB, trans)
		AssertError(t, err, tErr)
	})
}
