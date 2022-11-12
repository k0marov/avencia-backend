package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

//
// func TestTransferer(t *testing.T) {
// 	tRaw := RandomRawTransfer()
// 	transf := RandomTransfer()
// 	mockDB := NewStubDB()
// 	convert := func(values.RawTransfer) (values.Transfer, error) {
// 		return transf, nil
// 	}
// 	t.Run("error case - converting transfers throws", func(t *testing.T) {
// 		convert := func(gotTransf values.RawTransfer) (values.Transfer, error) {
// 			if gotTransf == tRaw {
// 				return values.Transfer{}, RandomError()
// 			}
// 			panic("unexpected")
// 		}
// 		err := service.NewTransferer(convert, nil, nil)(mockDB, tRaw)
// 		AssertSomeError(t, err)
// 	})
// 	validate := func(db.TDB, values.Transfer) error {
// 		return nil
// 	}
// 	t.Run("error case - validating transfers throws", func(t *testing.T) {
// 		err := RandomError()
// 		validate := func(gotDB db.TDB, transfer values.Transfer) error {
// 			if gotDB == mockDB && transfer == transf {
// 				return err
// 			}
// 			panic("unexpected")
// 		}
// 		gotErr := service.NewTransferer(convert, validate, nil)(mockDB, tRaw)
// 		AssertError(t, gotErr, err)
// 	})
//
// 	t.Run("happy case - forward to perform", func(t *testing.T) {
// 		err := RandomError()
// 		perform := func(gotDB db.TDB, t values.Transfer) error {
// 			if gotDB == mockDB && t == transf {
// 				return err
// 			}
// 			panic("unexpected")
// 		}
// 		gotErr := service.NewTransferer(convert, validate, perform)(mockDB, tRaw)
// 		AssertError(t, gotErr, err)
// 	})
// }
//
// func TestTransferPerformer(t *testing.T) {
// 	transf := values.Transfer{
// 		FromId: "John",
// 		ToId:   "Sam",
// 		SourceWalletId: "JohnWallet",
// 		Amount:   core.NewMoneyAmount(42),
// 	}
//
// 	withdrawTrans := transValues.Transaction{
// 		Source: transValues.TransSource{
// 			Type:   transValues.Transfer,
// 			Detail: "Sam",
// 		},
// 		WalletId: "JohnWallet",
// 		Money:   core.NewMoneyAmount(-42),
// 	}
// 	depositTrans := transValues.Transaction{
// 		Source: transValues.TransSource{
// 			Type:   transValues.Transfer,
// 			Detail: "John",
// 		},
// 		UserId: "Sam",
// 		Money: core.Money{
// 			Currency: "RUB",
// 			Amount:   core.NewMoneyAmount(42),
// 		},
// 	}
//
// 	mockDB := NewStubDB()
//
// 	t.Run("forward case", func(t *testing.T) {
// 		tErr := RandomError()
// 		transact := func(gotDB db.TDB, tList []transValues.Transaction) error {
// 			if gotDB == mockDB && reflect.DeepEqual(tList, []transValues.Transaction{withdrawTrans, depositTrans}) {
// 				return tErr
// 			}
// 			panic("unexpected")
// 		}
// 		gotErr := service.NewTransferPerformer(transact)(mockDB, transf)
// 		AssertError(t, gotErr, tErr)
// 	})
// }

func TestTransferConverter(t *testing.T) {
	mockDB := NewStubDB()
	rawTrans := RandomRawTransfer()
	user := RandomUser()
	wallet := RandomWallet()

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
		_, err := service.NewTransferConverter(userFromEmail, nil)(mockDB, rawTrans)
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - auth getter throws some other error", func(t *testing.T) {
		userFromEmail := func(string) (authEntities.User, error) {
			return authEntities.User{}, RandomError()
		}
		_, err := service.NewTransferConverter(userFromEmail, nil)(mockDB, rawTrans)
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
		_, err := service.NewTransferConverter(userFromEmail, getWallet)(mockDB, rawTrans)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		gotTrans, err := service.NewTransferConverter(userFromEmail, getWallet)(mockDB, rawTrans)
		AssertNoError(t, err)
		want := values.Transfer{
			FromId: rawTrans.FromId,
			ToId:   user.Id,
			SourceWallet: wallet,
			Amount:  rawTrans.Amount,
		}
		Assert(t, gotTrans, want, "converted transfers")
	})
}
