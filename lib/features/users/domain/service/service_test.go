package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	userEntities "github.com/AvenciaLab/avencia-backend/lib/features/users/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/users/domain/service"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestUserInfoGetter(t *testing.T) {
	userId := RandomString()
	wallet := RandomWallet()
	tLimits := RandomLimits()
	mockDB := NewStubDB()

	getWallet := func(db.DB, string) (walletEntities.Wallet, error) {
		return wallet, nil
	}
	t.Run("error case - getting wallets throws", func(t *testing.T) {
		getWallet := func(gotDB db.DB, user string) (walletEntities.Wallet, error) {
			if gotDB == mockDB && user  == userId {
				return walletEntities.Wallet{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet, nil)(mockDB, userId)
		AssertSomeError(t, err)
	})
	getLimits := func(db.DB, string) (limits.Limits, error) {
		return tLimits, nil
	}
	t.Run("error case - getting limits throws", func(t *testing.T) {
		getLimits := func(gotDB db.DB, user string) (limits.Limits, error) {
			if gotDB == mockDB && user == userId {
				return limits.Limits{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet, getLimits)(mockDB, userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		gotInfo, err := service.NewUserInfoGetter(getWallet, getLimits)(mockDB, userId)
		AssertNoError(t, err)
		Assert(t, gotInfo, userEntities.UserInfo{
			Id:     userId,
			Wallet: wallet,
			Limits: tLimits,
		}, "returned users info")
	})
}
