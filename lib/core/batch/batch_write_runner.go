package batch

import (
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
)

type WriteRunner = func(func(batch firestore_facade.WriteBatch) error) error

func NewWriteRunner(client firestore_facade.Simple) WriteRunner {
	return func(perform func(batch firestore_facade.WriteBatch) error) error {
		batch := client.Batch()
		err := perform(batch)
		if err != nil {
			return fmt.Errorf("performing (not committing) a batch write: %w", err)
		}
		_, err = batch.Commit(context.Background())
		if err != nil {
			return fmt.Errorf("committing a batch write: %w", err)
		}
		return nil
	}
}
