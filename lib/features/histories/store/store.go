package store

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/histories/store/mappers"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

func NewHistoryGetter(client *firestore.Client, decode mappers.TransEntriesDecoder) store.HistoryGetter {
	return func(userId string) ([]entities.TransEntry, error) {
		col := client.Collection("Users/" + userId + "/History")
		docs, err := col.Documents(context.Background()).GetAll()
		if err != nil {
			return []entities.TransEntry{}, core_err.Rethrow("fetching history entries from fs", err)
		}
		return decode(fs_facade.NewDocuments(docs))
	}
}

func NewTransStorer(getDoc fs_facade.DocGetter, encode mappers.TransEntryEncoder) store.TransStorer {
	return func(u fs_facade.Updater, t transValues.Transaction) error {
		doc := getDoc("Users/"+t.UserId+"/History/"+fs_facade.NewDocId())
		value := encode(t.Source, t.Money)
		u(doc, value)
		return nil 
	}
}
