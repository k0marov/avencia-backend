package store

import (
	"fmt"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/models"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/store/mappers"
)

func NewWithdrawsGetter(getDocs db.ColGetter, decode mappers.WithdrawsDecoder) store.WithdrawsGetter {
	return func(db db.DB, userId string) ([]models.Withdrawn, error) {
		path := fmt.Sprintf("Users/%s/Withdraws", userId)
		docs, err := getDocs(db, path)
		if err != nil {
			return []models.Withdrawn{}, fmt.Errorf("fetching a list of withdraws documents %w", err)
		}
		return decode(docs) 
	}

}

func NewWithdrawUpdater(updDoc db.Setter, encode mappers.WithdrawEncoder) store.WithdrawUpdater {
	return func(db db.DB, userId string, withdrawn core.Money) error {
		path := fmt.Sprintf("Users/%s/Withdraws/%s", userId, string(withdrawn.Currency))
		return updDoc(db, path, encode(withdrawn.Amount))
	}
}
