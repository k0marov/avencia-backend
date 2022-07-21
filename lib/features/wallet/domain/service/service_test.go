package service_test

import (
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
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
		wallet := map[string]float64{"USD": 400.0, "RUB": 42000.0, "BTC": 0.001}
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
