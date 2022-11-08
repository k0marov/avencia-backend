package store

import (

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/store"
)

func NewHistoryGetter(getDocs db.JsonColGetter[entities.TransEntry]) store.HistoryGetter {
	return func(db db.TDB, userId string) ([]entities.TransEntry, error) {
		col := []string{"histories", userId}
		return getDocs(db,col)
	}
}

func NewTransStorer(updDoc db.JsonSetter[entities.TransEntry]) store.EntryStorer {
	return func(db db.TDB, userId string, entry entities.TransEntry) error {
		path := []string{"histories", userId, general_helpers.RandomId()}
		return updDoc(db, path, entry)
	}
}
