package test_helpers

import (
	"errors"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
)

func NewStubDB() db.DB {
	return db.NewDB(newStubDB())
}

type stubDB struct {
	randId int
}

func newStubDB() stubDB {
	return stubDB{
		randId: RandomInt(),
	}
}

func (s stubDB) Get(path []string) (db.Document, error) {
	return db.Document{}, errors.New("unimplemented")
}
func (s stubDB) GetCollection(path []string) (db.Documents, error) {
	return db.Documents{}, errors.New("unimplemented")
}

func (s stubDB) Set(path []string, data []byte) error {
	return errors.New("unimplemented")
}

func (s stubDB) RunTransaction(func(db.DB) error) error {
	return errors.New("unimplemented")
}


