package foundationdb

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
)

type simpleDB struct {
	runTrans db.TransRunner
}

func NewSimpleDB(runTrans db.TransRunner) simpleDB {
  return simpleDB{runTrans: runTrans}
}

// Since in foundationdb everything must happen inside a transaction, 
// SimpleDB just forwards all calls to TransactionalDB, 
// except that it initiates and commits a separate transaction for every call

func (s simpleDB) Get(path []string) (db.Document, error) {
  var doc db.Document
  err := s.runTrans(func(db db.TDB) error {
    var err error
    doc, err = db.Get(path)
    return err
  })
  return doc, err
}
func (s simpleDB) GetCollection(path []string) (db.Documents, error) {
  var docs db.Documents
  err := s.runTrans(func(db db.TDB) error {
    var err error
    docs, err = db.GetCollection(path)
    return err
  })
  return docs, err
}
func (s simpleDB) Set(path []string, data []byte) error {
  return s.runTrans(func(db db.TDB) error {
    return db.Set(path, data)
  })
}
func (s simpleDB) Delete(path []string) error {
  return s.runTrans(func(db db.TDB) error {
    return db.Delete(path)
  })
}

