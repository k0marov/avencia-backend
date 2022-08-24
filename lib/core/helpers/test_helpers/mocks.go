package test_helpers

import "github.com/k0marov/avencia-backend/lib/core/db"

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

func (s stubDB) Get(path string) (db.Document, error) {
	return db.Document{}, nil
}

func (s stubDB) GetAll(path string) (db.Documents, error) {
	return db.Documents{}, nil
}

func (s stubDB) Set(path string, data map[string]any) error {
	return nil
}
