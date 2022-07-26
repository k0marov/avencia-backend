package batch

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
)

type WriteRunner = func(func(batch firestore_facade.BatchUpdater) error) error

// TODO: add a rethrow function which will just forward the error if it is a client error

func NewWriteRunner(client *firestore.Client) WriteRunner {
	return func(perform func(batch firestore_facade.BatchUpdater) error) error {
		batch := client.Batch()
		err := perform(firestore_facade.NewBatchUpdater(batch))
		if err != nil {
			return err
		}
		_, err = batch.Commit(context.Background())
		if err != nil {
			return fmt.Errorf("committing a batch write: %w", err)
		}
		return nil
	}
}
