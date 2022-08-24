package service_test

import (
	"reflect"
	"testing"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

func TestTransferer(t *testing.T) {
	tRaw := RandomRawTransfer()
	transf := RandomTransfer()
	mockDB := NewStubDB() 
	convert := func(values.RawTransfer) (values.Transfer, error) {
		return transf, nil
	}
	t.Run("error case - converting transfers throws", func(t *testing.T) {
		convert := func(gotTransf values.RawTransfer) (values.Transfer, error) {
			if gotTransf == tRaw {
				return values.Transfer{}, RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransferer(convert, nil, nil)(mockDB, tRaw)
		AssertSomeError(t, err)
	})
	validate := func(transfer values.Transfer) error {
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
		perform := func(gotDB db.DB, t values.Transfer) error {
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
		FromId: "John",
		ToId:   "Sam",
		Money: core.Money{
			Currency: "RUB",
			Amount:   core.NewMoneyAmount(42),
		},
	}

	withdrawTrans := transValues.Transaction{
		Source: transValues.TransSource{
			Type:   transValues.Transfer,
			Detail: "Sam",
		},
		UserId: "John",
		Money: core.Money{
			Currency: "RUB",
			Amount:   core.NewMoneyAmount(-42),
		},
	}
	depositTrans := transValues.Transaction{
		Source: transValues.TransSource{
			Type:   transValues.Transfer,
			Detail: "John",
		},
		UserId: "Sam",
		Money: core.Money{
			Currency: "RUB",
			Amount:   core.NewMoneyAmount(42),
		},
	}

	mockDB := NewStubDB()
	
	t.Run("forward case", func(t *testing.T) {
		tErr := RandomError()
		transact := func(gotDB db.DB, tList []transValues.Transaction) error {
			if gotDB == mockDB && reflect.DeepEqual(tList, []transValues.Transaction{withdrawTrans, depositTrans}) {
				return tErr
			}
			panic("unexpected")
		}
		gotErr := service.NewTransferPerformer(transact)(mockDB, transf)
		AssertError(t, gotErr, tErr)
	})
}

func TestTransferConverter(t *testing.T) {
	rawTrans := values.RawTransfer{
		FromId:  RandomString(),
		ToEmail: RandomString(),
		Money:   RandomPositiveMoney(),
	}
	user := RandomUser()

	userFromEmail := func(gotEmail string) (auth.User, error) {
		if gotEmail == rawTrans.ToEmail {
			return user, nil
		}
		panic("unexpected")
	}
	t.Run("error case - auth getter throws ErrNotFound", func(t *testing.T) {
		userFromEmail := func(string) (auth.User, error) {
			return auth.User{}, core_err.ErrNotFound
		}
		_, err := service.NewTransferConverter(userFromEmail)(rawTrans)
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - auth getter throws some other error", func(t *testing.T) {
		userFromEmail := func(string) (auth.User, error) {
			return auth.User{}, RandomError()
		}
		_, err := service.NewTransferConverter(userFromEmail)(rawTrans)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		gotTrans, err := service.NewTransferConverter(userFromEmail)(rawTrans)
		AssertNoError(t, err)
		want := values.Transfer{
			FromId: rawTrans.FromId,
			ToId:   user.Id,
			Money:  rawTrans.Money,
		}
		Assert(t, gotTrans, want, "converted transfers")
	})
}
