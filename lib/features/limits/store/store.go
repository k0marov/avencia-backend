package store

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/models"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/store/mappers"
)

func NewWithdrawsGetter(getDoc db.JsonGetter, decode mappers.WithdrawsDecoder) store.WithdrawsGetter {
	return func(db db.DB, userId string) (models.Withdraws, error) {
		path := []string{"withdrawn", userId}
		doc, err := getDoc(db, path)
		if err != nil && !core_err.IsNotFound(err) {
			return models.Withdraws{}, core_err.Rethrow("getting withdraws document", err)
		}
		return decode(doc)
	}

}

// TODO: consider adding more context info to every core_err.Rethrow()

func NewWithdrawUpdater(updDoc db.Setter, encode mappers.WithdrawEncoder) store.WithdrawUpdater {
	return func(db db.DB, userId string, withdrawn core.Money) error {
		path := []string{"withdrawn", userId}
		return updDoc(db, path, encode(withdrawn, time.Now()))
	}
}
