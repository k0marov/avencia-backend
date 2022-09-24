package foundationdb

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

type simpleDB struct {
  fdb fdb.Database
}

func NewSimpleDB(fdb fdb.Database) simpleDB {
  return simpleDB{fdb: fdb}
}

// Since in foundationdb everything must happen inside a transaction, 
// SimpleDB just forwards all calls to TransactionalDB, 
// except that it initiates and commits a separate transaction for every call

func (s simpleDB) Get(path []string) (db.Document, error) {
  var doc db.Document
  err := NewTransactionRunner(s.fdb)(func(dbHandle db.DB) error {
    var err error
    doc, err = db.GetterImpl(dbHandle, path)
    return err
  })
  return doc, err
}
func (s simpleDB) GetCollection(path []string) (db.Documents, error) {
  var docs db.Documents
  err := NewTransactionRunner(s.fdb)(func(dbHandle db.DB) error {
    var err error
    docs, err = db.ColGetterImpl(dbHandle, path)
    return err
  })
  return docs, err
}
func (s simpleDB) Set(path []string, data []byte) error {
  return NewTransactionRunner(s.fdb)(func(dbHandle db.DB) error {
    return db.SetterImpl(dbHandle, path, data)
  })
}
func (s simpleDB) Delete(path []string) error {
  return NewTransactionRunner(s.fdb)(func(dbHandle db.DB) error {
    return db.DeleterImpl(dbHandle, path)
  })
}

