package foundationdb

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/general_helpers"
)

type transactionalDB struct {
  t fdb.Transaction 
}
// FoundationDB interface is implemented both by fdb.Database and fdb.Transaction
type FoundationDB interface {
	Transact(func (fdb.Transaction) (interface{}, error)) (r interface{}, e error)
}

// NewTransactionRunner( fDB can be a fdb.Database instance)
func NewTransactionRunner(fDB FoundationDB) db.TransRunner {
	return func(perform func(db.TDB) error) error {
    _, err := fDB.Transact(func(t fdb.Transaction) (interface{}, error) {
      tDB := transactionalDB{t: t}
      err := perform(tDB)
      return nil, err
    })
    return err
	}
}

func pathToKey(path []string) fdb.Key {
  return fdb.Key(general_helpers.ConvTuple(path).Pack())
}


func (t transactionalDB) Get(path []string) (db.Document, error) {
  data, err := t.t.Get(pathToKey(path)).Get() 
  if err != nil {
    return db.Document{}, core_err.Rethrow("while getting a doc", err)
  }
  if data == nil {
  	return db.Document{}, core_err.ErrNotFound
  }
  return db.Document{
  	Path: path,
  	Data: &data,
  }, nil
}

func (t transactionalDB) GetCollection(path []string) (db.Documents, error) {
	res := t.t.GetRange(general_helpers.ConvTuple(path), fdb.RangeOptions{Mode: fdb.StreamingModeWantAll})
	kvs, err := res.GetSliceWithError()
	if err != nil {
		return db.Documents{}, core_err.Rethrow("getting slice of docs", err)
	}
	docs := db.Documents{} 
	for _, kv := range kvs {
		val := kv.Value
		docs = append(docs, db.Document{
			Path: path,
			Data: &val,
		})
	}
	return docs, nil
}

func (t transactionalDB) Set(path []string, data []byte) error {
  t.t.Set(pathToKey(path), []byte(data))
  return nil 
}

func (t transactionalDB) Delete(path []string) error {
	t.t.Clear(pathToKey(path)) 
	return nil
}

func (t transactionalDB) RunTransaction(perform func(db.TDB) error) error {
	_, err := t.t.Transact(func(trans fdb.Transaction) (interface{}, error) {
    err := perform(transactionalDB{
    	t: trans,
    })
    return nil, err
	})
	return err 
}

