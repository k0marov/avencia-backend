package validators_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestWalletOwnershipValidator(t *testing.T) {
	mockDB := NewStubDB()
	walletId := RandomString()
	callerId := RandomString()
	fittingWallet := wEntities.Wallet{
		Id: walletId,
		WalletVal: wEntities.WalletVal{
			OwnerId: callerId,
		},
	}
	notFittingWallet := wEntities.Wallet{
		Id: walletId,
		WalletVal: wEntities.WalletVal{
			OwnerId: RandomString(),
		},
	}

	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(gotDB db.TDB, id string) (wEntities.Wallet, error) {
			if gotDB == mockDB && id == walletId {
				return fittingWallet, RandomError()
			}
			panic("unexpected")
		}
		err := validators.NewWalletOwnershipValidator(getWallet)(mockDB, callerId, walletId)
		AssertSomeError(t, err)
	})
	t.Run("error case - the transaction initiator is not the owner of the provided wallet", func(t *testing.T) {
		getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
			return notFittingWallet, nil
		}
		err := validators.NewWalletOwnershipValidator(getWallet)(mockDB, callerId, walletId)
		AssertError(t, err, client_errors.Unauthorized)
	})
	t.Run("happy case", func(t *testing.T) {
		getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
			return fittingWallet, nil
		}
		err := validators.NewWalletOwnershipValidator(getWallet)(mockDB, callerId, walletId)
		AssertNoError(t, err)
	})
}

func TestTransactionValidator(t *testing.T) {
	mockDB := NewStubDB()
	tBalance := RandomPosMoneyAmount()
	trans := values.Transaction{
		Source:   RandomTransactionSource(),
		WalletId: RandomString(),
		Money:    RandomMoney(),
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

func TestWalletValidator(t *testing.T) {
	mockDB := NewStubDB()
	currency := RandomCurrency()
	trans := values.Transaction{
		Source:   RandomTransactionSource(),
		WalletId: RandomString(),
		Money:    core.Money{
			Currency: currency,
			Amount: core.NewMoneyAmount(-50.0),
		},
	}
	incorrectCurrencyWallet := wEntities.Wallet{
		WalletVal: wEntities.WalletVal{
			Currency: "IncorrectCurrency",
			Amount:   core.MoneyAmount{},
		},
	}
	notEnoughBalanceWallet := wEntities.Wallet{
		WalletVal: wEntities.WalletVal{
			Currency: currency,
			Amount: core.NewMoneyAmount(30),
		},
	}
	enoughBalanceWallet := wEntities.Wallet{
		WalletVal: wEntities.WalletVal{
			Currency: currency,
			Amount: core.NewMoneyAmount(100),
		},
	}
	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(gotDB db.TDB, walletId string) (wEntities.Wallet, error) {
			if gotDB == mockDB && walletId == trans.WalletId {
				return wEntities.Wallet{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := validators.NewWalletValidator(getWallet)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("error case - currencies don't match", func(t *testing.T) {
		getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
			return incorrectCurrencyWallet, nil 
		}
		_, err := validators.NewWalletValidator(getWallet)(mockDB, trans) 
		AssertError(t, err, client_errors.InvalidCurrency)
	})
	t.Run("error case - insufficient funds", func(t *testing.T) {
		getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
			return notEnoughBalanceWallet, nil
		}
		_, err := validators.NewWalletValidator(getWallet)(mockDB, trans)
		AssertError(t, err, client_errors.InsufficientFunds)
	})
	t.Run("happy case", func(t *testing.T) {
		getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
			return enoughBalanceWallet, nil
		}
		bal, err := validators.NewWalletValidator(getWallet)(mockDB, trans)
		AssertNoError(t, err)
		Assert(t, bal, enoughBalanceWallet.Amount, "returned current balance")
	})

}
