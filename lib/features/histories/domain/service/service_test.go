package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/service"
)

func TestHistoryGetter(t *testing.T) {
	userId := RandomString()
	mockDB := NewStubDB() 
	t.Run("error case - getting history entries from store throws", func(t *testing.T) {
		getFromStore := func(gotDB db.DB, gotUserId string) ([]entities.TransEntry, error) {
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
			CreatedAt: TimeWithYear(2022), 
		}
		entryOldest := entities.TransEntry{
			Source:    RandomTransactionSource(),
			Money:     RandomMoney(),
			CreatedAt: TimeWithYear(2000),
		}
		entryMiddle := entities.TransEntry{
			Source:    RandomTransactionSource(),
			Money:     RandomMoney(),
			CreatedAt: TimeWithYear(2010),
		}
		storeEntries := []entities.TransEntry{
			entryOldest, 
			entryNewest, 
			entryMiddle, 
		}
		getFromStore := func(gotDB db.DB, userId string) ([]entities.TransEntry, error) {
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

}
