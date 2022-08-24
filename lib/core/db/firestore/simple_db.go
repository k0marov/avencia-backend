package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/db"
)
type SimpleDB struct {
	c *firestore.Client
}
func (db SimpleDB) Get(path string) (db.Document, error) {
	doc, err := db.c.Doc(path).Get(context.Background())
	return newDocument(doc), err
}
func (db SimpleDB) GetAll(path string) (db.Documents, error) {
	docs, err := db.c.Collection(path).Documents(context.Background()).GetAll()
	return newDocuments(docs), err
}
func (db SimpleDB) Set(path string, data map[string]any) error {
	_, err := db.c.Doc(path).Set(context.Background(), data, firestore.MergeAll)
	return err
}
