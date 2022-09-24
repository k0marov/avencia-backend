package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
)

func TestBalanceGetter(t *testing.T) {
	userId := RandomString()
	mockDB := NewStubDB()
	t.Run("error case", func(t *testing.T) {
		getWallet := func(db.TDB, string) (entities.Wallet, error) {
			return entities.Wallet{}, RandomError()
		}
		_, err := service.NewBalanceGetter(getWallet)(mockDB, userId, RandomCurrency())
		AssertSomeError(t, err)
	})
	t.Run("should return 0 if there is no such currency in the wallet", func(t *testing.T) {
		wallet := entities.Wallet{}
		getWallet := func(gotDB db.TDB, user string) (entities.Wallet, error) {
			if gotDB == mockDB && user == userId {
				return wallet, nil
			}
			panic("unexpected")
		}
		balance, err := service.NewBalanceGetter(getWallet)(mockDB, userId, RandomCurrency())
		AssertNoError(t, err)
		Assert(t, balance, core.NewMoneyAmount(0), "returned balance")
	})
	t.Run("should return the value from wallets", func(t *testing.T) {
		wallet := entities.Wallet{"RUB": core.NewMoneyAmount(4000), "USD": core.NewMoneyAmount(300)}
		getWallet := func(db.TDB, string) (entities.Wallet, error) {
			return wallet, nil
		}
		balance, err := service.NewBalanceGetter(getWallet)(mockDB, userId, "USD")
		AssertNoError(t, err)
		Assert(t, balance, core.NewMoneyAmount(300), "returned balance")
	})
}
