package db

type DB struct {
	db dbInternal
}

func NewDB(db dbInternal) DB {
	return DB{
		db: db,
	}

}

type TransRunner = func(func(DB) error) error

type dbInternal interface {
	// Get should return core_err.NotFound if not found
	Get(path []string) (Document, error)
	GetCollection(path []string) (Documents, error)
	Set(path []string, data []byte) error
	Delete(path []string) error
	RunTransaction(func(DB) error) error
}

type Document struct {
	Path []string
	Data *[]byte 
}

type Documents []Document


type Deleter = func(db DB, path []string) error

func DeleterImpl(db DB, path []string) error {
	return db.db.Delete(path)
}



