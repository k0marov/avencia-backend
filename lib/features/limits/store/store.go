package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

// withdrawsDocGetter userId should be non-empty
type withdrawDocGetter = func(userId string, currency core.Currency) *firestore.DocumentRef

func NewWithdrawDocGetter(getDoc fs_facade.DocGetter) withdrawDocGetter {
	return func(userId string, currency core.Currency) *firestore.DocumentRef {
		doc := getDoc(fmt.Sprintf("Withdraws/%s/Withdraws/%s", userId, string(currency)))
		if doc == nil {
			panic("getting document ref for users's withdraws returned nil. Probably userId is empty.")
		}
		return doc
	}
}

const withdrawnKey = "withdrawn"

// TODO: somehow simplify this

func NewWithdrawsGetter(client *firestore.Client) store.WithdrawsGetter {
	return func(userId string) (map[string]values.WithdrawnWithUpdated, error) {
		coll := client.Collection(fmt.Sprintf("Withdraws/%s/Withdraws", userId))
		docs, err := coll.DocumentRefs(context.Background()).GetAll()
		if err != nil {
			return map[string]values.WithdrawnWithUpdated{}, fmt.Errorf("fetching a list of withdraws documents %w", err)
		}

		withdraws := map[string]values.WithdrawnWithUpdated{}

		for _, doc := range docs {
			snap, err := doc.Get(context.Background())
			if err != nil {
				return map[string]values.WithdrawnWithUpdated{}, core_err.Rethrow("fetching a withdraw document", err)
			}
			withdrawnVal := snap.Data()[withdrawnKey]
			withdrawn, err := general_helpers.DecodeFloat(withdrawnVal)
			if err != nil {
				return map[string]values.WithdrawnWithUpdated{}, err
			}
			withdraws[doc.ID] = values.WithdrawnWithUpdated{
				Withdrawn: core.NewMoneyAmount(withdrawn),
				UpdatedAt: snap.UpdateTime,
			}
		}
		return withdraws, nil
	}
}

func NewWithdrawUpdater(getDoc withdrawDocGetter) store.WithdrawUpdater {
	return func(update fs_facade.Updater, userId string, withdrawn core.Money) error {
		doc := getDoc(userId, withdrawn.Currency)
		return update(doc, map[string]any{withdrawnKey: withdrawn.Amount.Num()})
	}
}
