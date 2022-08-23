package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/db"
)

type DBTransaction = *firestore.Transaction

type TransactionalDB struct {
	t DBTransaction
	c *firestore.Client
}


type TransactionDBFactory = func(t DBTransaction) TransactionalDB 

func NewTransactionDBFactory(c *firestore.Client) TransactionDBFactory {
	return func(t DBTransaction) TransactionalDB {
		return TransactionalDB{
			t: t,
			c: c,
		}
	}
}


func (db TransactionalDB) get(path string) (db.Document, error) {
	doc, err := db.t.Get(db.c.Doc(path))
	return newDocument(doc), err
}

func (db TransactionalDB) getAll(path string) (db.Documents, error) {
	docs, err := db.t.Documents(db.c.Collection(path)).GetAll()
	return newDocuments(docs), err
}

func (db TransactionalDB) set(path string, data map[string]any) error {
	return db.t.Set(db.c.Doc(path), data, firestore.MergeAll)
}
