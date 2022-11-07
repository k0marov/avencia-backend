package validators_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func TestTransactionValidator(t *testing.T) {
	mockDB := NewStubDB()
	tBalance := RandomPosMoneyAmount()
	trans := values.Transaction{
		Source: RandomTransactionSource(),
		WalletId: RandomString(),
		Money:   core.NewMoneyAmount(50.0),
	}
	checkLimit := func(db.TDB, values.Transaction) error {
		return nil
	}
	t.Run("error case - limit checker throws", func(t *testing.T) {
		err := RandomError()
		checkLimit := func(gotDB db.TDB, t values.Transaction) error {
			if gotDB == mockDB && t == trans {
				return err
			}
			panic("unexpected")
		}
		_, gotErr := validators.NewTransactionValidator(checkLimit, nil)(mockDB, trans)
		AssertError(t, gotErr, err)
	})
	t.Run("forward case - forward to enoughBalanceValidator", func(t *testing.T) {
		tErr := RandomError()
    enoughBalanceValidator := func(gotDB db.TDB, gotTrans values.Transaction) (core.MoneyAmount, error) {
    	if gotDB == mockDB && gotTrans == trans {
    		return tBalance, tErr
    	}
    	panic("unexpected")
    }
    gotBalance, gotErr := validators.NewTransactionValidator(checkLimit, enoughBalanceValidator)(mockDB, trans)
    AssertError(t, gotErr, tErr)
    Assert(t, gotBalance, tBalance, "returned balance")
	})
}

func TestEnoughBalanceValidator(t *testing.T) {
	mockDB := NewStubDB()
	trans := values.Transaction{
		Source:   RandomTransactionSource(),
		WalletId: RandomString(),
		Money:    core.NewMoneyAmount(-50.0),
	}
	notEnoughBalance := core.NewMoneyAmount(30)
	enoughBalance := core.NewMoneyAmount(100)
	t.Run("error case - getting balance throws", func(t *testing.T) {
		getBalance := func(gotDB db.TDB, walletId string) (core.MoneyAmount, error) {
			if gotDB == mockDB && walletId == trans.WalletId {
				return core.MoneyAmount{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := validators.NewEnoughBalanceValidator(getBalance)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("error case - insufficient funds", func(t *testing.T) {
		getBalance := func(db.TDB, string) (core.MoneyAmount, error) {
			return notEnoughBalance, nil
		}
		_, err := validators.NewEnoughBalanceValidator(getBalance)(mockDB, trans)
		AssertError(t, err, client_errors.InsufficientFunds)
	})
	t.Run("happy case", func(t *testing.T) {
		getBalance := func(db.TDB, string) (core.MoneyAmount, error) {
			return enoughBalance, nil
		}
		bal, err := validators.NewEnoughBalanceValidator(getBalance)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, bal, enoughBalance, "returned current balance")
	})

}
