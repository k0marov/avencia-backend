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

// TODO: add checking if Money.Amount is negative to Transferer, otherwise everyone will be able to transfer money to their account from other users :)

func TestTransferConverter(t *testing.T) {
	rawTrans := values.RawTransfer{
		FromId:  RandomString(),
		ToEmail: RandomString(),
		Money:   RandomMoney(),
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
