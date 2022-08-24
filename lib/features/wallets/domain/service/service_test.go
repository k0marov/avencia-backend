package service_test

import (
	"testing"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
)

func TestWalletGetter(t *testing.T) {
	userId := RandomString()
	mockDB := NewStubDB() 
	t.Run("error case - store throws", func(t *testing.T) {
		getWallet := func(db.DB, string) (map[string]any, error) {
			return nil, RandomError()
		}
		_, err := service.NewWalletGetter(getWallet)(mockDB, userId)
		AssertSomeError(t, err)
	})
	t.Run("error case - some value is not a float", func(t *testing.T) {
		getWallet := func(db.DB, string) (map[string]any, error) {
			return map[string]any{"test": "not-a-float"}, nil
		}
		_, err := service.NewWalletGetter(getWallet)(mockDB, userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		storedWallet := map[string]any{"USD": 400.0, "RUB": 42000.0, "BTC": 0.001}
		wallet := map[core.Currency]core.MoneyAmount{"USD": core.NewMoneyAmount(400.0), "RUB": core.NewMoneyAmount(42000.0), "BTC": core.NewMoneyAmount(0.001)}
		getWallet := func(gotDB db.DB, user string) (map[string]any, error) {
			if gotDB == mockDB && user == userId {
				return storedWallet, nil
			}
			panic("unexpected")
		}
		gotWallet, err := service.NewWalletGetter(getWallet)(mockDB, userId)
		AssertNoError(t, err)
		Assert(t, gotWallet, wallet, "returned wallets")
	})
}

func TestBalanceGetter(t *testing.T) {
	userId := RandomString()
	mockDB := NewStubDB()
	t.Run("error case", func(t *testing.T) {
		getWallet := func(db.DB, string) (entities.Wallet, error) {
			return entities.Wallet{}, RandomError()
		}
		_, err := service.NewBalanceGetter(getWallet)(mockDB, userId, RandomCurrency())
		AssertSomeError(t, err)
	})
	t.Run("should return 0 if there is no such currency in wallets", func(t *testing.T) {
		wallet := entities.Wallet{}
		getWallet := func(gotDB db.DB, user string) (entities.Wallet, error) {
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
		getWallet := func(db.DB, string) (entities.Wallet, error) {
			return wallet, nil
		}
		balance, err := service.NewBalanceGetter(getWallet)(mockDB, userId, "USD")
		AssertNoError(t, err)
		Assert(t, balance, core.NewMoneyAmount(300), "returned balance")
	})
}
