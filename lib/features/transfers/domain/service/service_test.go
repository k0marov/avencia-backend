package service_test

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
	"reflect"
	"testing"
)

func TestTransferer(t *testing.T) {
	tRaw := RandomRawTransfer()
	transf := RandomTransfer()
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
		err := service.NewTransferer(convert, nil, nil, nil)(tRaw)
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
		gotErr := service.NewTransferer(convert, validate, nil, nil)(tRaw)
		AssertError(t, gotErr, err)
	})

	// TODO: simplify this callback hell
	runBatch := func(f func(firestore_facade.BatchUpdater) error) error {
		return f(func(*firestore.DocumentRef, map[string]any) error {
			return nil
		})
	}

	withdrawTrans := transValues.Transaction{
		UserId: transf.FromId,
		Money: core.Money{
			Currency: transf.Money.Currency,
			Amount:   transf.Money.Amount.Neg(),
		},
	}
	depositTrans := transValues.Transaction{
		UserId: transf.ToId,
		Money: core.Money{
			Currency: transf.Money.Currency,
			Amount:   transf.Money.Amount,
		},
	}

	transact := func(firestore_facade.BatchUpdater, transValues.Transaction) error { return nil }

	t.Run("error case - withdrawing from caller fails", func(t *testing.T) {
		transact := func(u firestore_facade.BatchUpdater, t transValues.Transaction) error {
			if reflect.DeepEqual(t, withdrawTrans) {
				return RandomError()
			}
			return nil
		}
		err := service.NewTransferer(convert, validate, runBatch, transact)(tRaw)
		AssertSomeError(t, err)
	})
	t.Run("error case - depositing to recipient fails", func(t *testing.T) {
		transact := func(u firestore_facade.BatchUpdater, t transValues.Transaction) error {
			if t == depositTrans {
				return RandomError()
			}
			return nil
		}
		err := service.NewTransferer(convert, validate, runBatch, transact)(tRaw)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		err := service.NewTransferer(convert, validate, runBatch, transact)(tRaw)
		AssertNoError(t, err)
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

func TestTransferValidator(t *testing.T) {
	t.Run("error case - money.amount is negative", func(t *testing.T) {
		trans := values.Transfer{
			FromId: RandomString(),
			ToId:   RandomString(),
			Money:  RandomNegativeMoney(),
		}
		err := service.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.NegativeTransferAmount)
	})
	t.Run("error case - transfering to yourself", func(t *testing.T) {
		user := RandomId()
		trans := values.Transfer{
			FromId: user,
			ToId:   user,
			Money:  RandomPositiveMoney(),
		}
		err := service.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.TransferingToYourself)
	})
	t.Run("error case - transfering 0", func(t *testing.T) {
		trans := values.Transfer{
			FromId: RandomString(),
			ToId:   RandomString(),
			Money: core.Money{
				Currency: RandomCurrency(),
				Amount:   core.NewMoneyAmount(0),
			},
		}
		err := service.NewTransferValidator()(trans)
		AssertError(t, err, client_errors.TransferingZero)
	})
}
