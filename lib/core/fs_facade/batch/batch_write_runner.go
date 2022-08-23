package batch

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
)

type WriteRunner = func(func(batch fs_facade.BatchUpdater) error) error

func NewWriteRunner(client *firestore.Client) WriteRunner {
	return func(perform func(batch fs_facade.BatchUpdater) error) error {
		batch := client.RunTransaction(context.Background(), func(ctx context.Context, t *firestore.Transaction) error {

			return nil
		})
		err := perform(fs_facade.NewBatchUpdater(batch))
		if err != nil {
			return err
		}
		_, err = batch.Commit(context.Background())
		if err != nil {
			return core_err.Rethrow("committing a batch write", err)
		}
		return nil
	}
}
