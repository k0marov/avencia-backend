package service_test

import (
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	"testing"
)

func TestWalletGetter(t *testing.T) {
	userId := RandomString()
	t.Run("error case - store throws", func(t *testing.T) {
		getWallet := func(string) (map[string]any, error) {
			return nil, RandomError()
		}
		_, err := service.NewWalletGetter(getWallet)(userId)
		AssertSomeError(t, err)
	})
	t.Run("error case - some value is not a float", func(t *testing.T) {
		getWallet := func(string) (map[string]any, error) {
			return map[string]any{"test": "not-a-float"}, nil
		}
		_, err := service.NewWalletGetter(getWallet)(userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		storedWallet := map[string]any{"USD": 400.0, "RUB": 42000.0, "BTC": 0.001}
		wallet := map[core.Currency]core.MoneyAmount{"USD": 400.0, "RUB": 42000.0, "BTC": 0.001}
		getWallet := func(user string) (map[string]any, error) {
			if user == userId {
				return storedWallet, nil
			}
			panic("unexpected")
		}
		gotWallet, err := service.NewWalletGetter(getWallet)(userId)
		AssertNoError(t, err)
		Assert(t, gotWallet, wallet, "returned wallet")
	})
}

func TestBalanceGetter(t *testing.T) {
	userId := RandomString()
	t.Run("error case", func(t *testing.T) {
		getWallet := func(user string) (entities.Wallet, error) {
			return entities.Wallet{}, RandomError()
		}
		_, err := service.NewBalanceGetter(getWallet)(userId, RandomCurrency())
		AssertSomeError(t, err)
	})
	t.Run("should return 0 if there is no such currency in wallet", func(t *testing.T) {
		wallet := entities.Wallet{}
		getWallet := func(user string) (entities.Wallet, error) {
			if user == userId {
				return wallet, nil
			}
			panic("unexpected")
		}
		balance, err := service.NewBalanceGetter(getWallet)(userId, RandomCurrency())
		AssertNoError(t, err)
		Assert(t, balance, 0, "returned balance")
	})
	t.Run("should return the value from wallet", func(t *testing.T) {
		wallet := entities.Wallet{"RUB": 4000, "USD": 300}
		getWallet := func(string) (entities.Wallet, error) {
			return wallet, nil
		}
		balance, err := service.NewBalanceGetter(getWallet)(userId, "USD")
		AssertNoError(t, err)
		Assert(t, balance, core.MoneyAmount(300), "returned balance")
	})
}
