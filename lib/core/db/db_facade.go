package db

import (
	"time"
)

type DB interface {
  get(path string) (Document, error)
  getAll(path string) (Documents, error) 
  set(path string, data map[string]any) error
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
  return db.get(path) 
}
func ColGetterImpl(db DB, path string) (Documents, error) {
	return db.getAll(path)
}
func SetterImpl(db DB, path string, data map[string]any) error {
  return db.set(path, data)
}


