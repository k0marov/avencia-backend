package service_test

import (
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	userEntities "github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
	"testing"
)

func TestUserInfoGetter(t *testing.T) {
	userId := RandomString()
	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(user string) (walletEntities.Wallet, error) {
			if user == userId {
				return walletEntities.Wallet{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet)(userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		wallet := walletEntities.Wallet{"RUB": 1000, "USD": 100.5}
		getWallet := func(string) (walletEntities.Wallet, error) {
			return wallet, nil
		}
		gotInfo, err := service.NewUserInfoGetter(getWallet)(userId)
		AssertNoError(t, err)
		Assert(t, gotInfo, userEntities.UserInfo{
			Id:     userId,
			Wallet: wallet,
		}, "returned user info")
	})
}
