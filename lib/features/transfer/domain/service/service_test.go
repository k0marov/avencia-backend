package service_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_errors"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
	"testing"
)

// TODO: add checking if Money.Amount is negative in Transferer, otherwise everyone will be able to transfer money to their account from other users :)

func TestTransferer(t *testing.T) {
	tRaw := RandomRawTransfer()
	//transf := RandomTransfer()
	t.Run("error case - money.amount is negative", func(t *testing.T) {
		tRaw := values.RawTransfer{
			FromId:  RandomString(),
			ToEmail: RandomString(),
			Money:   RandomNegativeMoney(),
		}
		err := service.NewTransferer(nil, nil)(tRaw)
		AssertError(t, err, client_errors.NegativeTransferAmount)
	})
	t.Run("error case - converting transfer throws", func(t *testing.T) {
		convert := func(gotTransf values.RawTransfer) (values.Transfer, error) {
			if gotTransf == tRaw {
				return values.Transfer{}, RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransferer(convert, nil)(tRaw)
		AssertSomeError(t, err)
	})
	t.Run("making withdraw from caller transaction fails", func(t *testing.T) {

	})
	t.Run("making deposit to recepient transaction fails", func(t *testing.T) {

	})
	t.Run("happy case", func(t *testing.T) {

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
