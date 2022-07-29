package service_test

import (
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	limitsEntities "github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	userEntities "github.com/k0marov/avencia-backend/lib/features/users/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/users/domain/service"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallets/domain/entities"
	"testing"
)

func TestUserInfoGetter(t *testing.T) {
	userId := RandomString()
	wallet := RandomWallet()
	limits := RandomLimits()

	getWallet := func(string) (walletEntities.Wallet, error) {
		return wallet, nil
	}
	t.Run("error case - getting wallets throws", func(t *testing.T) {
		getWallet := func(user string) (walletEntities.Wallet, error) {
			if user == userId {
				return walletEntities.Wallet{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet, nil)(userId)
		AssertSomeError(t, err)
	})
	getLimits := func(string) (limitsEntities.Limits, error) {
		return limits, nil
	}
	t.Run("error case - getting limits throws", func(t *testing.T) {
		getLimits := func(user string) (limitsEntities.Limits, error) {
			if user == userId {
				return limitsEntities.Limits{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet, getLimits)(userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		gotInfo, err := service.NewUserInfoGetter(getWallet, getLimits)(userId)
		AssertNoError(t, err)
		Assert(t, gotInfo, userEntities.UserInfo{
			Id:     userId,
			Wallet: wallet,
			Limits: limits,
		}, "returned users info")
	})
}
