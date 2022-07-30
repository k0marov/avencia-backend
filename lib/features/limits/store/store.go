package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/limits/store/mappers"
)

func NewWithdrawsGetter(client *firestore.Client, decode mappers.WithdrawsDecoder) store.WithdrawsGetter {
	return func(userId string) ([]values.WithdrawnModel, error) {
		col := client.Collection(fmt.Sprintf("Users/%s/Withdraws", userId))
		docs, err := col.Documents(context.Background()).GetAll()
		if err != nil {
			return []values.WithdrawnModel{}, fmt.Errorf("fetching a list of withdraws documents %w", err)
		}
		
		return decode(fs_facade.NewDocuments(docs)) 
	}

}

func NewWithdrawUpdater(getDoc fs_facade.DocGetter, encode mappers.WithdrawEncoder) store.WithdrawUpdater {
	return func(update fs_facade.Updater, userId string, withdrawn core.Money) error {
		doc := getDoc(fmt.Sprintf("Users/%s/Withdraws/%s", userId, string(withdrawn.Currency)))
		return update(doc, encode(withdrawn.Amount))
	}
}
