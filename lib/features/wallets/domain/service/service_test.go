package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
)

func TestWalletCreator(t *testing.T) {
	mockDB := NewStubDB()
	userId := RandomString()
	currency := RandomCurrency()
	data := service.WalletCreationData{
		UserId:   userId,
		Currency: currency,
	}

	t.Run("forward case", func(t *testing.T) {
		tId := RandomString()
		tErr := RandomError()
		storeCreator := func(gotDB db.TDB, wallet entities.WalletVal) (string, error) {
			if gotDB == mockDB &&
				wallet.OwnerId == userId &&
				wallet.Currency == currency && wallet.Amount.IsEqual(core.NewMoneyAmount(0)) {
				return tId, tErr
			}
			panic("unexpected")
		}
		id, err := service.NewWalletCreator(storeCreator)(mockDB, data)
		AssertError(t, err, tErr)
		Assert(t, id, tId, "returned id")
	})

}
