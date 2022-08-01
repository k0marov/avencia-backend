package fs_facade

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// TODO: prettify


func NewDocId() string {
	uuid, _ := uuid.NewUUID() 
	return uuid.String()
}

type Document struct {
	Id        string
	Data      map[string]any
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Documents []Document

func NewDocument(doc *firestore.DocumentSnapshot) Document {
	return Document{
		Id:        doc.Ref.ID,
		Data:      doc.Data(),
		UpdatedAt: doc.UpdateTime,
		CreatedAt: doc.CreateTime,
	}
}

func NewDocuments(docs []*firestore.DocumentSnapshot) (res Documents) {
	for _, doc := range docs {
		res = append(res, NewDocument(doc))
	}
	return
}

// TODO: add tests for the store layers thanks to the new facades

type DocGetter = func(path string) *firestore.DocumentRef

func NewDocGetter(client *firestore.Client) DocGetter {
	return func(path string) *firestore.DocumentRef {
		return client.Doc(path)
	}
}

type Updater = func(doc *firestore.DocumentRef, data map[string]any) error

// BatchUpdater is used to show that a method wants to execute its actions only in a batch
type BatchUpdater Updater

func NewBatchUpdater(batch *firestore.WriteBatch) BatchUpdater {
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
