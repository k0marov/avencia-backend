package validators_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
	"testing"
)

func TestTransferValidator(t *testing.T) {
	t.Run("error case - money.amount is negative", func(t *testing.T) {
		trans := values.Transfer{
			FromId: RandomString(),
			ToId:   RandomString(),
			Money:  RandomNegativeMoney(),
		}
		err := validators.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.NegativeTransferAmount)
	})
	t.Run("error case - transfering to yourself", func(t *testing.T) {
		user := RandomId()
		trans := values.Transfer{
			FromId: user,
			ToId:   user,
			Money:  RandomPositiveMoney(),
		}
		err := validators.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.TransferringToYourself)
	})
	t.Run("error case - transferring 0", func(t *testing.T) {
		trans := values.Transfer{
			FromId: RandomString(),
			ToId:   RandomString(),
			Money: core.Money{
				Currency: RandomCurrency(),
				Amount:   core.NewMoneyAmount(0),
			},
		}
		err := validators.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.TransferringZero)
	})
}
