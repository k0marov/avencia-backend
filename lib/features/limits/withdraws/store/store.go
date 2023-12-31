package store

import (
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/store"
)

func NewWithdrawsGetter(getDoc db.JsonGetter[models.Withdraws]) store.WithdrawsGetter {
	return func(db db.TDB, userId string) (models.Withdraws, error) {
		path := []string{"withdrawn", userId}
		withdraws, err := getDoc(db, path)
		if core_err.IsNotFound(err) {
			return models.Withdraws{}, nil
		} 
		if err != nil {
			return models.Withdraws{}, core_err.Rethrow("getting withdraws doc", err)
		}
		return withdraws, nil
	}

}


func NewWithdrawUpdater(updDoc db.JsonUpdater[models.WithdrawVal]) store.WithdrawUpdater {
	return func(db db.TDB, userId string, withdrawn core.Money) error {
		path := []string{"withdrawn", userId}
		val := models.WithdrawVal{
			Withdrawn: withdrawn.Amount,
			UpdatedAt: time.Now().Unix(), // TODO: this should in the business logic
		}
		return updDoc(db, path, string(withdrawn.Currency), val)
	}
}
