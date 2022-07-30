package validators_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	"testing"
)

func TestTransactionValidator(t *testing.T) {
	curBalance := core.NewMoneyAmount(100.0)
	trans := values.Transaction{
		Source: RandomTransactionSource(),
		UserId: RandomString(),
		Money: core.Money{
			Currency: RandomCurrency(),
			Amount:   core.NewMoneyAmount(50.0),
		},
	}
	checkLimit := func(t values.Transaction) error {
		return nil
	}
	t.Run("error case - limit checker throws", func(t *testing.T) {
		err := RandomError()
		checkLimit := func(t values.Transaction) error {
			if t == trans {
				return err
			}
			panic("unexpected")
		}
		_, gotErr := validators.NewTransactionValidator(checkLimit, nil)(trans)
		AssertError(t, gotErr, err)
	})
	t.Run("error case - getting balance throws", func(t *testing.T) {
		getBalance := func(string, core.Currency) (core.MoneyAmount, error) {
			return core.NewMoneyAmount(0), RandomError()
		}
		_, err := validators.NewTransactionValidator(checkLimit, getBalance)(trans)
		AssertSomeError(t, err)
	})
	t.Run("error case - insufficient funds", func(t *testing.T) {
		getBalance := func(string, core.Currency) (core.MoneyAmount, error) {
			return core.NewMoneyAmount(30.0), nil
		}
		trans := values.Transaction{
			Money: core.Money{
				Amount: core.NewMoneyAmount(-50.0),
			},
		}
		_, err := validators.NewTransactionValidator(checkLimit, getBalance)(trans)
		AssertError(t, err, client_errors.InsufficientFunds)
	})
	t.Run("happy case", func(t *testing.T) {
		getBalance := func(userId string, currency core.Currency) (core.MoneyAmount, error) {
			if userId == trans.UserId && currency == trans.Money.Currency {
				return curBalance, nil
			}
			panic("unexpected")
		}
		bal, err := validators.NewTransactionValidator(checkLimit, getBalance)(trans)
		AssertNoError(t, err)
		Assert(t, bal, curBalance, "returned current balance")
	})
}
