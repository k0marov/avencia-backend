package store

import (
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/store"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func NewHistoryGetter(getDocs db.JsonColGetter[entities.TransEntry]) store.HistoryGetter {
	return func(db db.DB, userId string) ([]entities.TransEntry, error) {
		col := []string{"histories", userId}
		return getDocs(db,col)
	}
}

func NewTransStorer(updDoc db.JsonSetter[entities.TransEntry]) store.TransStorer {
	return func(db db.DB, t transValues.Transaction) error {
		path := []string{"histories", t.UserId, general_helpers.RandomId()}
		tEntry := entities.TransEntry{
			Source:    t.Source,
			Money:     t.Money,
			CreatedAt: time.Now(), // TODO: this should be in the business layer
		}
		return updDoc(db, path, tEntry)
	}
}
