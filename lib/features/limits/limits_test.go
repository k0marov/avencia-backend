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
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestLimitChecker(t *testing.T) {
	tLimit := limits.Limit{
		Withdrawn: core.NewMoneyAmount(500),
		Max:       core.NewMoneyAmount(600),
	}
	wallet := RandomString()
	mockDB := NewStubDB()

	getLimit := func(gotDB db.TDB, walletId string) (limits.Limit, error) {
		if gotDB == mockDB && wallet == walletId {
			return tLimit, nil
		}
		panic("unexpected")
	}

	t.Run("error case - getting limit throws", func(t *testing.T) {
		getLimit := func(db.TDB, string) (limits.Limit, error) {
			return limits.Limit{}, RandomError()
		}
		err := limits.NewLimitChecker(getLimit)(mockDB, transValues.Transaction{Money: RandomNegativeMoney()})
		AssertSomeError(t, err)
	})
	t.Run("error case - limit exceeded", func(t *testing.T) {
		trans := transValues.Transaction{
			Source:   RandomTransactionSource(),
			WalletId: wallet,
			Money: core.Money{
				Amount: core.NewMoneyAmount(-200),
			},
		}
		err := limits.NewLimitChecker(getLimit)(mockDB, trans)
		AssertError(t, err, client_errors.WithdrawLimitExceeded)
	})
	t.Run("happy case", func(t *testing.T) {
		trans := transValues.Transaction{
			Source:   RandomTransactionSource(),
			WalletId: wallet,
			Money: core.Money{
				Amount: core.NewMoneyAmount(-50),
			},
		}
		err := limits.NewLimitChecker(getLimit)(mockDB, trans)
		AssertNoError(t, err)
	})
	t.Run("happy case - value is positive (its a deposit)", func(t *testing.T) {
		trans := transValues.Transaction{
			Source:   RandomTransactionSource(),
			WalletId: wallet,
			Money: core.Money{
				Amount: core.NewMoneyAmount(1000),
			},
		}
		err := limits.NewLimitChecker(getLimit)(mockDB, trans)
		AssertNoError(t, err)
	})
}

func TestLimitGetter(t *testing.T) {
	user := RandomString()
	mockDB := NewStubDB()
	tLimits := limits.Limits{
		"RUB": RandomLimit(),
		"USD": RandomLimit(),
	}
	wallet := wEntities.Wallet{
		Id: RandomString(),
		WalletVal: wEntities.WalletVal{
			OwnerId:  user,
			Currency: "RUB",
			Amount:   RandomMoneyAmount(),
		},
	}

	getWallet := func(gotDB db.TDB, wId string) (wEntities.Wallet, error) {
		if gotDB == mockDB && wId == wallet.Id {
			return wallet, nil
		}
		panic("unexpected")
	}
	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
			return wEntities.Wallet{}, RandomError()
		}
		_, err := limits.NewLimitGetter(getWallet, nil)(mockDB, wallet.Id)
		AssertSomeError(t, err)
	})

	getLimits := func(db.TDB, string) (limits.Limits, error) {
		return tLimits, nil
	}
	t.Run("error case - getting limits throws", func(t *testing.T) {
		getLimits := func(gotDB db.TDB, userId string) (limits.Limits, error) {
			if gotDB == mockDB && userId == user {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		_, err := limits.NewLimitGetter(getWallet, getLimits)(mockDB, wallet.Id)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		limit, err := limits.NewLimitGetter(getWallet, getLimits)(mockDB, wallet.Id)
		AssertNoError(t, err)
		Assert(t, limit, tLimits["RUB"], "the returned limit")
	})
}

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
