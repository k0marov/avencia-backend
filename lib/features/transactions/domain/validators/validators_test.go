package validators_test

import (
	"testing"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

func TestTransactionValidator(t *testing.T) {
	mockDB := NewStubDB() 
	curBalance := core.NewMoneyAmount(100.0)
	trans := values.Transaction{
		Source: RandomTransactionSource(),
		UserId: RandomString(),
		Money: core.Money{
			Currency: RandomCurrency(),
			Amount:   core.NewMoneyAmount(50.0),
		},
	}
	checkLimit := func(db.DB, values.Transaction) error {
		return nil
	}
	t.Run("error case - limit checker throws", func(t *testing.T) {
		err := RandomError()
		checkLimit := func(gotDB db.DB, t values.Transaction) error {
			if gotDB == mockDB && t == trans {
				return err
			}
			panic("unexpected")
		}
		_, gotErr := validators.NewTransactionValidator(checkLimit, nil)(mockDB, trans)
		AssertError(t, gotErr, err)
	})
	t.Run("error case - getting balance throws", func(t *testing.T) {
		getBalance := func(db.DB, string, core.Currency) (core.MoneyAmount, error) {
			return core.NewMoneyAmount(0), RandomError()
		}
		_, err := validators.NewTransactionValidator(checkLimit, getBalance)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("error case - insufficient funds", func(t *testing.T) {
		getBalance := func(db.DB, string, core.Currency) (core.MoneyAmount, error) {
			return core.NewMoneyAmount(30.0), nil
		}
		trans := values.Transaction{
			Money: core.Money{
				Amount: core.NewMoneyAmount(-50.0),
			},
		}
		_, err := validators.NewTransactionValidator(checkLimit, getBalance)(mockDB, trans)
		AssertError(t, err, client_errors.InsufficientFunds)
	})
	t.Run("happy case", func(t *testing.T) {
		getBalance := func(gotDB db.DB, userId string, currency core.Currency) (core.MoneyAmount, error) {
			if gotDB == mockDB && userId == trans.UserId && currency == trans.Money.Currency {
				return curBalance, nil
			}
			panic("unexpected")
		}
		bal, err := validators.NewTransactionValidator(checkLimit, getBalance)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, bal, curBalance, "returned current balance")
	})
}
