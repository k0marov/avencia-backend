package store

import (
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/histories/store/mappers"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

func NewHistoryGetter(getDocs db.JsonCollectionGetter, decode mappers.TransEntriesDecoder) store.HistoryGetter {
	return func(db db.DB, userId string) ([]entities.TransEntry, error) {
		col := []string{"histories", userId}
		docs, err := getDocs(db, col)
		if err != nil {
			return []entities.TransEntry{}, core_err.Rethrow("fetching history entries from fs", err)
		}
		return decode(docs)
	}
}

func NewTransStorer(updDoc db.JsonSetter, encode mappers.TransEntryEncoder) store.TransStorer {
	return func(db db.DB, t transValues.Transaction) error {
		path := []string{"histories", t.UserId, general_helpers.RandomId()}
		value := encode(t.Source, t.Money)
		return updDoc(db, path, value)
	}
}
