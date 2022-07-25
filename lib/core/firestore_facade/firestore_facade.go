package firestore_facade

import (
	"cloud.google.com/go/firestore"
	"context"
)

// TODO: add tests for the store layers thanks to the new facades

type DocGetter = func(path string) *firestore.DocumentRef

func NewDocGetter(client *firestore.Client) DocGetter {
	return func(path string) *firestore.DocumentRef {
		return client.Doc(path)
	}
}

type Updater = func(doc *firestore.DocumentRef, data map[string]any) error

func NewBatchUpdater(batch *firestore.WriteBatch) Updater {
	return func(doc *firestore.DocumentRef, data map[string]any) error {
		batch.Set(doc, data, firestore.MergeAll)
		return nil
	}
}

func NewSimpleUpdater() Updater {
	return func(doc *firestore.DocumentRef, data map[string]any) error {
		_, err := doc.Set(context.Background(), data, firestore.MergeAll)
		return err
	}
}
