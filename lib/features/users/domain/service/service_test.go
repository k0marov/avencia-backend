package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	"github.com/AvenciaLab/avencia-backend/lib/features/users/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/users/domain/service"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestUserInfoGetter(t *testing.T) {
	detUser := RandomDetailedUser()
	userId := RandomString()
	wallets := []walletEntities.Wallet{RandomWallet(), RandomWallet()}
	tLimits := RandomLimits()
	mockDB := NewStubDB()

	getWallet := func(db.TDB, string) ([]walletEntities.Wallet, error) {
		return wallets, nil
	}
	t.Run("error case - getting wallets throws", func(t *testing.T) {
		getWallets := func(gotDB db.TDB, user string) ([]walletEntities.Wallet, error) {
			if gotDB == mockDB && user == userId {
				return []walletEntities.Wallet{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallets, nil, nil)(mockDB, userId)
		AssertSomeError(t, err)
	})
	getLimits := func(db.TDB, string) (limits.Limits, error) {
		return tLimits, nil
	}
	t.Run("error case - getting limits throws", func(t *testing.T) {
		getLimits := func(gotDB db.TDB, user string) (limits.Limits, error) {
			if gotDB == mockDB && user == userId {
				return limits.Limits{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet, getLimits, nil)(mockDB, userId)
		AssertSomeError(t, err)
	})

	getUser := func(string) (authEntities.DetailedUser, error) {
		return detUser, nil 
	}
	t.Run("error case - getting detailed user info throws", func(t *testing.T) {
		getUser := func(user string) (authEntities.DetailedUser, error) {
			if user == userId {
				return authEntities.DetailedUser{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet, getLimits, getUser)(mockDB, userId) 
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		gotInfo, err := service.NewUserInfoGetter(getWallet, getLimits, getUser)(mockDB, userId)
		AssertNoError(t, err)
		Assert(t, gotInfo, entities.UserInfo{
			User:   detUser,
			Wallets: wallets,
			Limits: tLimits,
		}, "returned users info")
	})
}
