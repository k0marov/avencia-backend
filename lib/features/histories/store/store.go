package store

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/histories/store/mappers"
)

func NewHistoryGetter(client *firestore.Client, decode mappers.TransEntriesDecoder) store.HistoryGetter {
  return func(userId string) ([]entities.TransEntry, error) {
    col := client.Collection("Users/"+userId+"/History")
    docs, err := col.Documents(context.Background()).GetAll()
    if err != nil {
      return []entities.TransEntry{}, core_err.Rethrow("fetching history entries from fs", err)
    }
    return decode(fs_facade.NewDocuments(docs))
  }
}

