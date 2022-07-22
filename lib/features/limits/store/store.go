package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

// withdrawsDocGetter userId should be non-empty
type withdrawDocGetter = func(userId string, currency core.Currency) *firestore.DocumentRef

func NewWithdrawsDocGetter(client firestore_facade.Simple) withdrawDocGetter {
	return func(userId string, currency core.Currency) *firestore.DocumentRef {
		doc := client.Doc(fmt.Sprintf("Withdraws/%s/Withdraws/%s", userId, string(currency)))
		if doc == nil {
			panic("getting document ref for user's withdraws returned nil. Probably userId is empty.")
		}
		return doc
	}
}

func NewWithdrawsGetter(client *firestore.Client) store.WithdrawsGetter {
	return func(userId string) (res map[string]values.WithdrawnWithUpdated, err error) {
		coll := client.Collection(fmt.Sprintf("Withdraws/%s/Withdraws", userId))
		docs, err := coll.DocumentRefs(context.Background()).GetAll()
		if err != nil {
			return map[string]values.WithdrawnWithUpdated{}, fmt.Errorf("fetching a list of withdraws documents %w", err)
		}
		for _, doc := range docs {
			snap, err := doc.Get(context.Background())
			if err != nil {
				return map[string]values.WithdrawnWithUpdated{}, fmt.Errorf("fetching a withdraw document: %w", err)
			}
			withdrawnVal := snap.Data()["withdrawn"]
			withdrawn, ok := withdrawnVal.(float64)
			if !ok {
				return map[string]values.WithdrawnWithUpdated{}, fmt.Errorf("withdrawn value %v is not a float", withdrawnVal)
			}
			res[doc.ID] = values.WithdrawnWithUpdated{
				Withdrawn: core.MoneyAmount(withdrawn),
				UpdatedAt: snap.UpdateTime,
			}
		}
		return
	}
}
