package db

type Document struct {
	Path []string
	Data *[]byte 
}
type Documents []Document

// SDB is a simple db, meaning that every operations runs as a non-transactional call to the db.
type SDB interface {
	// Get should return core_err.NotFound if not found
	Get(path []string) (Document, error)
	GetCollection(path []string) (Documents, error)
	Set(path []string, data []byte) error
	Delete(path []string) error
}

// TDB is a transactional db, meaning that every call invoked on an instance of this handle
// runs in 1 transaction, which will later be commited to the db.
type TDB interface {
	SDB
	RunTransaction(func(TDB) error) error
}
type TransRunner = func(func(TDB) error) error

type Getter = func(db SDB, path []string) (Document, error) 
type ColGetter = func(db SDB, path []string) (Documents, error) 
type Setter = func(db SDB, path []string, data []byte) error 
type Deleter = func(db SDB, path []string) error



