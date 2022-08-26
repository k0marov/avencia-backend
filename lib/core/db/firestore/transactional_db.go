package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/db"
)

type TransactionalDB struct {
	t *firestore.Transaction
	c *firestore.Client
}


// TODO: it turns out, in firestore transactions all reads must happen before writes. Rewrite everything 

func NewTransactionRunner(c *firestore.Client ) db.TransactionRunner {
	return func(perform func(db.DB) error) error {
		return c.RunTransaction(context.Background(), func(ctx context.Context, t *firestore.Transaction) error {
			db := db.NewDB(TransactionalDB{
				t: t, 
				c: c, 
			})
			return perform(db)
		})
	}
}




func (db TransactionalDB) Get(path string) (db.Document, error) {
	doc, err := db.t.Get(db.c.Doc(path))
	return newDocument(doc), err
}

func (db TransactionalDB) GetAll(path string) (db.Documents, error) {
	docs, err := db.t.Documents(db.c.Collection(path)).GetAll()
	return newDocuments(docs), err
}

func (db TransactionalDB) Set(path string, data map[string]any) error {
	return db.t.Set(db.c.Doc(path), data, firestore.MergeAll)
}
