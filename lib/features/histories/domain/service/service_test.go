package service_test

import (
	"testing"

	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
)

func TestHistoryGetter(t *testing.T) {
	userId := RandomString()
	t.Run("error case - getting history entries from store throws", func(t *testing.T) {
		getFromStore := func(gotUserId string) ([]entities.TransEntry, error) {
			if gotUserId == userId {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewHistoryGetter(getFromStore)(userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case - should returned entries from store sorted by createdAt", func(t *testing.T) {
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
		getFromStore := func(userId string) ([]entities.TransEntry, error) {
			return storeEntries, nil
		}

		wantEntries := []entities.TransEntry{
			entryNewest, 
			entryMiddle, 
			entryOldest, 
		}

		gotEntries, err := service.NewHistoryGetter(getFromStore)(userId) 
		AssertNoError(t, err)
		Assert(t, gotEntries, wantEntries, "sorted entries")
	})
}

func TestTransStorer(t *testing.T) {

}
