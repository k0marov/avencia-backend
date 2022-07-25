package service_test

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_errors"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
	"testing"
)

func TestTransferer(t *testing.T) {
	tRaw := RandomRawTransfer()
	transf := RandomTransfer()
	t.Run("error case - money.amount is negative", func(t *testing.T) {
		tRaw := values.RawTransfer{
			FromId:  RandomString(),
			ToEmail: RandomString(),
			Money:   RandomNegativeMoney(),
		}
		err := service.NewTransferer(nil, nil, nil)(tRaw)
		AssertError(t, err, client_errors.NegativeTransferAmount)
	})
	convert := func(values.RawTransfer) (values.Transfer, error) {
		return transf, nil
	}
	t.Run("error case - converting transfer throws", func(t *testing.T) {
		convert := func(gotTransf values.RawTransfer) (values.Transfer, error) {
			if gotTransf == tRaw {
				return values.Transfer{}, RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransferer(convert, nil, nil)(tRaw)
		AssertSomeError(t, err)
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
			Amount:   -transf.Money.Amount,
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
			if t == withdrawTrans {
				return RandomError()
			}
			return nil
		}
		err := service.NewTransferer(convert, runBatch, transact)(tRaw)
		AssertSomeError(t, err)
	})
	t.Run("error case - depositing to recipient fails", func(t *testing.T) {
		transact := func(u firestore_facade.BatchUpdater, t transValues.Transaction) error {
			if t == depositTrans {
				return RandomError()
			}
			return nil
		}
		err := service.NewTransferer(convert, runBatch, transact)(tRaw)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		err := service.NewTransferer(convert, runBatch, transact)(tRaw)
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
			return auth.User{}, core_errors.ErrNotFound
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
		Assert(t, gotTrans, want, "converted transfer")
	})
}
