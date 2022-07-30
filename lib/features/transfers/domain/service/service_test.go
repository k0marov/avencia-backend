package service_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
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
		err := service.NewTransferer(convert, nil, nil)(tRaw)
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
		gotErr := service.NewTransferer(convert, validate, nil)(tRaw)
		AssertError(t, gotErr, err)
	})

	t.Run("happy case - forward to perform", func(t *testing.T) {
		err := RandomError()
		perform := func(t values.Transfer) error {
			if t == transf {
				return err
			}
			panic("unexpected")
		}
		gotErr := service.NewTransferer(convert, validate, perform)(tRaw)
		AssertError(t, gotErr, err)
	})
}

func TestTransferPerformer(t *testing.T) {
	transf := RandomTransfer()

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

	transact := func(fs_facade.BatchUpdater, transValues.Transaction) error { return nil }

	t.Run("error case - withdrawing from caller fails", func(t *testing.T) {
		transact := func(u fs_facade.BatchUpdater, t transValues.Transaction) error {
			if reflect.DeepEqual(t, withdrawTrans) {
				return RandomError()
			}
			return nil
		}
		err := service.NewTransferPerformer(StubRunBatch, transact)(transf)
		AssertSomeError(t, err)
	})
	t.Run("error case - depositing to recipient fails", func(t *testing.T) {
		transact := func(u fs_facade.BatchUpdater, t transValues.Transaction) error {
			if t == depositTrans {
				return RandomError()
			}
			return nil
		}
		err := service.NewTransferPerformer(StubRunBatch, transact)(transf)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		err := service.NewTransferPerformer(StubRunBatch, transact)(transf)
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
