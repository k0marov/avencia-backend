package db

import (
	"time"
)


type DB struct {
	db dbInternal
}
func NewDB(db dbInternal) DB {
	return DB{
		db: db,
	}

}


type TransactionRunner = func(func(DB) error) error



type dbInternal interface {
  Get(path string) (Document, error)
  GetAll(path string) (Documents, error) 
  Set(path string, data map[string]any) error
}

type Document struct {
	Id        string
	Data      map[string]any
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Documents []Document

type DocGetter = func(db DB, path string) (Document, error)
type ColGetter = func(db DB, path string) (Documents, error) 
type Setter = func(db DB, path string, data map[string]any) error

func DocGetterImpl(db DB, path string) (Document, error) {
  return db.db.Get(path) 
}
func ColGetterImpl(db DB, path string) (Documents, error) {
	return db.db.GetAll(path)
}
func SetterImpl(db DB, path string, data map[string]any) error {
  return db.db.Set(path, data)
}

