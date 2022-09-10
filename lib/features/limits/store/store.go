package store

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/models"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
)

func NewWithdrawsGetter(getDoc db.JsonGetter[models.Withdraws]) store.WithdrawsGetter {
	return func(db db.DB, userId string) (models.Withdraws, error) {
		path := []string{"withdrawn", userId}
		return getDoc(db, path)
	}

}

// TODO: consider adding more context info to every core_err.Rethrow()

func NewWithdrawUpdater(updDoc db.JsonUpdater[models.WithdrawVal]) store.WithdrawUpdater {
	return func(db db.DB, userId string, withdrawn core.Money) error {
		path := []string{"withdrawn", userId}
		val := models.WithdrawVal{
			Withdrawn: withdrawn.Amount,
			UpdatedAt: time.Now(), // TODO: this should in the business logic
		}
		return updDoc(db, path, string(withdrawn.Currency), val)
	}
}
