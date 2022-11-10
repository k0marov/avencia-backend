package validators_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestWalletOwnershipValidator(t *testing.T) {
	mockDB := NewStubDB()
	trans := RandomMetaTrans()
	fittingWallet := walletEntities.Wallet{
		WalletVal: walletEntities.WalletVal{
			OwnerId: trans.CallerId,
		},
	}
	notFittingWallet := walletEntities.Wallet{
		WalletVal: walletEntities.WalletVal{
			OwnerId: RandomString(),
		},
	}

	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(gotDB db.TDB, id string) (walletEntities.Wallet, error) {
			if gotDB == mockDB && id == trans.WalletId {
				return fittingWallet, RandomError()
			}
			panic("unexpected")
		}
		err := validators.NewWalletOwnershipValidator(getWallet)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("error case - the transaction initiator is not the owner of the provided wallet", func(t *testing.T) {
		getWallet := func(db.TDB, string) (walletEntities.Wallet, error) {
			return notFittingWallet, nil
		}
		err := validators.NewWalletOwnershipValidator(getWallet)(mockDB, trans)
		AssertError(t, err, client_errors.Unauthorized)
	})
	t.Run("happy case", func(t *testing.T) {
		getWallet := func(db.TDB, string) (walletEntities.Wallet, error) {
			return fittingWallet, nil
		}
		err := validators.NewWalletOwnershipValidator(getWallet)(mockDB, trans)
		AssertNoError(t, err)
	})
}

func TestTransactionValidator(t *testing.T) {
	mockDB := NewStubDB()
	tBalance := RandomPosMoneyAmount()
	trans := values.Transaction{
		Source:   RandomTransactionSource(),
		WalletId: RandomString(),
		Money:    core.NewMoneyAmount(50.0),
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
	notEnoughBalanceWallet := walletEntities.Wallet{
		WalletVal: walletEntities.WalletVal{
			Amount: core.NewMoneyAmount(30),
		},
	}
	enoughBalanceWallet := walletEntities.Wallet{
		WalletVal: walletEntities.WalletVal{
			Amount: core.NewMoneyAmount(100),
		},
	}
	t.Run("error case - getting balance throws", func(t *testing.T) {
		getWallet := func(gotDB db.TDB, walletId string) (walletEntities.Wallet, error) {
			if gotDB == mockDB && walletId == trans.WalletId {
				return walletEntities.Wallet{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := validators.NewEnoughBalanceValidator(getWallet)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("error case - insufficient funds", func(t *testing.T) {
		getWallet := func(db.TDB, string) (walletEntities.Wallet, error) {
			return notEnoughBalanceWallet, nil
		}
		_, err := validators.NewEnoughBalanceValidator(getWallet)(mockDB, trans)
		AssertError(t, err, client_errors.InsufficientFunds)
	})
	t.Run("happy case", func(t *testing.T) {
		getWallet := func(db.TDB, string) (walletEntities.Wallet, error) {
			return enoughBalanceWallet, nil
		}
		bal, err := validators.NewEnoughBalanceValidator(getWallet)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, bal, enoughBalanceWallet.Amount, "returned current balance")
	})

}
