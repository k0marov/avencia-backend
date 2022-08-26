package foundationdb

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
)

type transactionalDB struct {
  t fdb.Transaction 
}

func NewTransactionRunner(fDB fdb.Transaction) db.TransactionRunner {
	return func(perform func(db.DB) error) error {
    _, err := fDB.Transact(func(t fdb.Transaction) (interface{}, error) {
      tDB := transactionalDB{t: t}
      err := perform(db.NewDB(tDB))
      return nil, err
    })
    return err
	}
}

func pathToKey(path []string) fdb.Key {
  return fdb.Key(tuple.Tuple{path}.Pack())
}

func (t transactionalDB) Get(path []string) (db.Document, error) {
  res := t.t.Get(pathToKey(path))
  data, err := res.Get()
  if err != nil {
    return db.Document{}, core_err.Rethrow("while getting a doc", err)
  }
  return db.Document{
  	Path: path,
  	Data: data,
  }, nil
}

func (t transactionalDB) GetCollection(path []string) (db.Documents, error) {
  panic("unimplemented")
}

func (t transactionalDB) Set(path []string, data []byte) error {
  t.t.Set(pathToKey(path), []byte(data))
  return nil 
}
