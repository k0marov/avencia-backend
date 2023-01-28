package validators_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestTransferValidator(t *testing.T) {
	caller := RandomString()
	wallet := wEntities.Wallet{
		Id:        "",
		WalletVal: wEntities.WalletVal{OwnerId: caller},
	}
	t.Run("error case - money.amount is negative", func(t *testing.T) {
		trans := values.Transfer{
			Money: RandomNegativeMoney(),
			FromWallet: wallet,
		}
		err := validators.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.NegativeTransferAmount)
	})
	t.Run("error case - transfering to yourself", func(t *testing.T) {
		trans := values.Transfer{
			FromId: caller,
			ToId:   caller,
			FromWallet: wallet,
			Money: RandomPositiveMoney(),
		}
		err := validators.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.TransferringToYourself)
	})
	t.Run("error case - transferring 0", func(t *testing.T) {
		trans := values.Transfer{
			FromId: caller,
			FromWallet: wallet,
			Money: core.Money{Amount: core.NewMoneyAmount(0)},
		}
		err := validators.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.TransferringZero)
	})
	t.Run("error case - caller is not the wallet owner", func(t *testing.T) {
		trans := values.Transfer{
			FromId: RandomString(), 
			FromWallet: wallet,
			Money: RandomPositiveMoney(),
		}
		err := validators.NewTransferValidator()(trans) 
		AssertError(t, err, client_errors.Unauthorized)
	})
	t.Run("happy case", func(t *testing.T) {
		trans := values.Transfer{
			FromId:caller,
			ToId:         RandomString(),
			FromWallet: wallet,
			Money:       RandomPositiveMoney(),
		}
		err := validators.NewTransferValidator()(trans) 
		AssertNoError(t, err)
	})
}
