package crud

import (
	"fmt"

	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
)

type Entity interface {
  Id() string
}

type Creator[E Entity] func(e E) error
type Updater[E Entity] func(e E) error
type Reader[E Entity] func(id string) (E, error)
type Deleter[E Entity] func(id string) error

type CRUDStore[E Entity] struct {
	db db.DB
	pathPrefix []string
}

func NewCRUDStore[E Entity](db db.DB, pathPrefix []string) CRUDStore[E] {
  return CRUDStore[E]{db: db, pathPrefix: pathPrefix}
}

func (s CRUDStore[E]) getEntityPath(id string) []string {
  return append(s.pathPrefix, id)
}

func (s CRUDStore[E]) Create(e E) error {
  path := s.getEntityPath(e.Id())
  err := db.JsonSetterImpl(s.db, path, e) 
  if err != nil {
    return core_err.Rethrow(fmt.Sprintf("while creating a CRUD entity %+v at path %+v", e, path), err)
  }
  return nil
}

func (s CRUDStore[E]) Update(e E) error {
  path := s.getEntityPath(e.Id())
  err := db.JsonMultiUpdaterImpl(s.db, path, e)
  if err != nil {
    return core_err.Rethrow(fmt.Sprintf("while updating a CRUD entity %+v at path %+v", e, path), err)
  }
  return nil
}

func (s CRUDStore[E]) Read(id string) (E, error) {
  path := s.getEntityPath(id) 
  e, err := db.JsonGetterImpl[E](s.db, path)
  if err != nil {
    return e, core_err.Rethrow(fmt.Sprintf("while getting a CRUD entity at path %+v", path), err)
  }
  return e, nil
}

func (s CRUDStore[E]) Delete(id string) error {
  path := s.getEntityPath(id) 
  err := db.DeleterImpl(s.db, path)
  if err != nil {
    return core_err.Rethrow(fmt.Sprintf("while deleting a CRUD entity at path %+v", path), err)
  }
  return nil
}
