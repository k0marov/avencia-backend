package service_test

import (
	"reflect"
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestTransferer(t *testing.T) {
	tRaw := RandomRawTransfer()
	transf := RandomTransfer()
	mockDB := NewStubDB()
	convert := func(db.TDB, values.RawTransfer) (values.Transfer, error) {
		return transf, nil
	}
	t.Run("error case - converting transfers throws", func(t *testing.T) {
		convert := func(gotDB db.TDB, gotTransf values.RawTransfer) (values.Transfer, error) {
			if gotDB == mockDB && gotTransf == tRaw {
				return values.Transfer{}, RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransferer(convert, nil, nil)(mockDB, tRaw)
		AssertSomeError(t, err)
	})
	validate := func(values.Transfer) error {
		return nil
	}
	t.Run("error case - validating transfers throws", func(t *testing.T) {
		err := RandomError()
		validate := func(transfer values.Transfer) error {
			if transfer == transf {
				return err
			}
			panic("unexpected")
		}
		gotErr := service.NewTransferer(convert, validate, nil)(mockDB, tRaw)
		AssertError(t, gotErr, err)
	})

	t.Run("happy case - forward to perform", func(t *testing.T) {
		err := RandomError()
		perform := func(gotDB db.TDB, t values.Transfer) error {
			if gotDB == mockDB && t == transf {
				return err
			}
			panic("unexpected")
		}
		gotErr := service.NewTransferer(convert, validate, perform)(mockDB, tRaw)
		AssertError(t, gotErr, err)
	})
}

func TestTransferPerformer(t *testing.T) {
	transf := values.Transfer{
		FromId:     "John",
		ToId:       "Sam",
		FromWallet: wEntities.Wallet{Id: "JohnWallet"},
		ToWallet:   wEntities.Wallet{Id: "SamWallet"},
		Money: core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(42),
		},
	}

	wantWithdrawTrans := transValues.Transaction{
		Source: transValues.TransSource{
			Type:   transValues.Transfer,
			Detail: "Sam",
		},
		WalletId: "JohnWallet",
		Money: core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(-42),
		},
	}
	wantDepositTrans := transValues.Transaction{
		Source: transValues.TransSource{
			Type:   transValues.Transfer,
			Detail: "John",
		},
		WalletId: "SamWallet",
		Money: core.Money{
			Currency: "USD",
			Amount:   core.NewMoneyAmount(42),
		},
	}

	mockDB := NewStubDB()

	t.Run("forward case", func(t *testing.T) {
		tErr := RandomError()
		transact := func(gotDB db.TDB, tList []transValues.Transaction) error {
			if gotDB == mockDB && reflect.DeepEqual(tList, []transValues.Transaction{wantWithdrawTrans, wantDepositTrans}) {
				return tErr
			}
			panic("unexpected")
		}
		gotErr := service.NewTransferPerformer(transact)(mockDB, transf)
		AssertError(t, gotErr, tErr)
	})
}

func TestWalletFinder(t *testing.T) {
	mockDB := NewStubDB()
	curr := RandomCurrency()
	wallets := []wEntities.Wallet{
		RandomWallet(),
		RandomWallet(),
		{WalletVal: wEntities.WalletVal{Currency: curr}},
		RandomWallet(),
	}
	userId := RandomString()

	getWallets := func(gotDB db.TDB, user string) ([]wEntities.Wallet, error) {
		if gotDB == mockDB && user == userId {
			return wallets, nil
		}
		panic("unexpected")
	}

	t.Run("error case - getting wallets throws", func(t *testing.T) {
		getWallets := func(db.TDB, string) ([]wEntities.Wallet, error) {
			return wallets, RandomError()
		}
		_, err := service.NewWalletFinder(getWallets)(mockDB, userId, curr)
		AssertSomeError(t, err)
	})
	t.Run("error case - proper wallet is not found", func(t *testing.T) {
		_, err := service.NewWalletFinder(getWallets)(mockDB, userId, core.Currency(RandomString()))
		AssertError(t, err, client_errors.ProperWalletNotFound)
	})
	t.Run("happy case", func(t *testing.T) {
		wallet, err := service.NewWalletFinder(getWallets)(mockDB, userId, curr)
		AssertNoError(t, err)
		Assert(t, wallet, wallets[2], "returned wallet")
	})
}

func TestTransferConverter(t *testing.T) {
	mockDB := NewStubDB()
	rawTrans := RandomRawTransfer()
	user := RandomUser()
	wallet := RandomWallet()
	toWallet := RandomWallet()

	userFromEmail := func(gotEmail string) (authEntities.User, error) {
		if gotEmail == rawTrans.ToEmail {
			return user, nil
		}
		panic("unexpected")
	}
	t.Run("error case - auth getter throws ErrNotFound", func(t *testing.T) {
		userFromEmail := func(string) (authEntities.User, error) {
			return authEntities.User{}, core_err.ErrNotFound
		}
		_, err := service.NewTransferConverter(userFromEmail, nil, nil)(mockDB, rawTrans)
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - auth getter throws some other error", func(t *testing.T) {
		userFromEmail := func(string) (authEntities.User, error) {
			return authEntities.User{}, RandomError()
		}
		_, err := service.NewTransferConverter(userFromEmail, nil, nil)(mockDB, rawTrans)
		AssertSomeError(t, err)
	})
	getWallet := func(gotDB db.TDB, walletId string) (wEntities.Wallet, error) {
		if gotDB == mockDB && walletId == rawTrans.SourceWalletId {
			return wallet, nil
		}
		panic("unexpected")
	}
	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
			return wEntities.Wallet{}, RandomError()
		}
		_, err := service.NewTransferConverter(userFromEmail, getWallet, nil)(mockDB, rawTrans)
		AssertSomeError(t, err)
	})
	findWallet := func(gotDB db.TDB, userId string, curr core.Currency) (wEntities.Wallet, error) {
		if gotDB == mockDB && userId == user.Id && curr == wallet.Currency {
			return toWallet, nil
		}
		panic("unexpected")
	}
	t.Run("error case - wallet finder throws", func(t *testing.T) {
		findWallet := func(db.TDB, string, core.Currency) (wEntities.Wallet, error) {
			return wEntities.Wallet{}, RandomError()
		}
		_, err := service.NewTransferConverter(userFromEmail, getWallet, findWallet)(mockDB, rawTrans)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		gotTrans, err := service.NewTransferConverter(userFromEmail, getWallet, findWallet)(mockDB, rawTrans)
		AssertNoError(t, err)
		want := values.Transfer{
			FromId:     rawTrans.FromId,
			ToId:       user.Id,
			FromWallet: wallet,
			ToWallet:   toWallet,
			Money: core.Money{
				Currency: wallet.Currency,
				Amount:   rawTrans.Amount,
			},
		}
		Assert(t, gotTrans, want, "converted transfers")
	})
}
