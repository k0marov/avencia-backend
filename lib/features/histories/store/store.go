package store

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/store"
)

func NewHistoryGetter(client *firestore.Client) store.HistoryGetter {
  return func(userId string) (fs_facade.Documents, error) {
    col := client.Collection("Users/"+userId+"/History")
    docs, err := col.Documents(context.Background()).GetAll()
    if err != nil {
      return fs_facade.Documents{}, core_err.Rethrow("fetching history entries from fs", err)
    }
    return fs_facade.NewDocuments(docs), nil
  }
}

