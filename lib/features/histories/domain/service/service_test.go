package service_test

import (
	"testing"
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/service"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func TestHistoryGetter(t *testing.T) {
	userId := RandomString()
	mockDB := NewStubDB() 
	t.Run("error case - getting history entries from store throws", func(t *testing.T) {
		getFromStore := func(gotDB db.TDB, gotUserId string) ([]entities.TransEntry, error) {
			if gotDB == mockDB && gotUserId == userId {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewHistoryGetter(getFromStore)(mockDB, userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case - should return entries from store sorted by createdAt", func(t *testing.T) {
		entryNewest := entities.TransEntry{
			Source:    RandomTransactionSource(),
			Money:     RandomMoney(),
			CreatedAt: TimeWithYear(2022).Unix(), 
		}
		entryOldest := entities.TransEntry{
			Source:    RandomTransactionSource(),
			Money:     RandomMoney(),
			CreatedAt: TimeWithYear(2000).Unix(),
		}
		entryMiddle := entities.TransEntry{
			Source:    RandomTransactionSource(),
			Money:     RandomMoney(),
			CreatedAt: TimeWithYear(2010).Unix(),
		}
		storeEntries := []entities.TransEntry{
			entryOldest, 
			entryNewest, 
			entryMiddle, 
		}
		getFromStore := func(gotDB db.TDB, userId string) ([]entities.TransEntry, error) {
			return storeEntries, nil
		}

		wantEntries := []entities.TransEntry{
			entryNewest, 
			entryMiddle, 
			entryOldest, 
		}

		gotEntries, err := service.NewHistoryGetter(getFromStore)(mockDB, userId) 
		AssertNoError(t, err)
		Assert(t, gotEntries, wantEntries, "sorted entries")
	})
}

func TestTransStorer(t *testing.T) {
	mockDB := NewStubDB()
	wallet := RandomWalletInfo()
	trans := RandomTransactionData()	
	getWallet := func(gotDB db.TDB, walletId string) (walletEntities.WalletInfo, error) {
		if gotDB == mockDB && walletId == trans.WalletId {
			return wallet, nil
		}
		panic("unexpected")
	}
	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(db.TDB, string) (walletEntities.WalletInfo, error) {
      return walletEntities.WalletInfo{}, RandomError()
		}
		err := service.NewEntryStorer(getWallet, nil)(mockDB, trans)
		AssertSomeError(t, err)
	})
	t.Run("forward case - forward to store", func(t *testing.T) {
		tErr := RandomError()
		storeEntry := func(gotDB db.TDB, userId string, entry entities.TransEntry) error {
			if gotDB == mockDB && userId == wallet.OwnerId && 
				entry.Money.Currency == wallet.Money.Currency && 
				entry.Money.Amount == trans.Money && entry.Source == trans.Source && 
				TimeAlmostEqual(time.Unix(entry.CreatedAt, 0), time.Now()) {
        	return tErr
				} 
			panic("unexpected")
		}
		err := service.NewEntryStorer(getWallet, storeEntry)(mockDB, trans)
		AssertError(t, err, tErr)
	})
}
