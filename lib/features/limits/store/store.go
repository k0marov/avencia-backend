package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
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

func NewWithdrawsGetter(client *firestore.Client) store.WithdrawsGetter {
	return func(userId string) (fs_facade.Documents, error) {
		col := client.Collection(fmt.Sprintf("Withdraws/%s/Withdraws", userId))
		docs, err := col.Documents(context.Background()).GetAll()
		if err != nil {
			return fs_facade.Documents{}, core_err.Rethrow("fetching all withdraws docs from fs", err)
		}
		return fs_facade.NewDocuments(docs), nil
	}
}

func NewWithdrawUpdater(getDoc withdrawDocGetter) store.WithdrawUpdater {
	return func(update fs_facade.Updater, userId string, withdrawn core.Money) error {
		doc := getDoc(userId, withdrawn.Currency)
		return update(doc, map[string]any{withdrawnKey: withdrawn.Amount.Num()})
	}
}
