package batch

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
)

type WriteRunner = func(func(batch firestore_facade.BatchUpdater) error) error

// TODO: refactor the core directory

func NewWriteRunner(client *firestore.Client) WriteRunner {
	return func(perform func(batch firestore_facade.BatchUpdater) error) error {
		batch := client.Batch()
		err := perform(firestore_facade.NewBatchUpdater(batch))
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
