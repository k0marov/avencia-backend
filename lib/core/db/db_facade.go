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

type Getter = func(db DB, path []string) (Document, error) 
type ColGetter = func(db DB, path []string) (Documents, error) 
type Setter = func(db DB, path []string, data []byte) error 
type Deleter = func(db DB, path []string) error

func GetterImpl(db DB, path []string) (Document, error) {
	return db.db.Get(path) 
}
func ColGetterImpl(db DB, path []string) (Documents, error) {
	return db.db.GetCollection(path)
}
func SetterImpl(db DB, path []string, data []byte) error {
	return db.db.Set(path, data)
}
func DeleterImpl(db DB, path []string) error {
	return db.db.Delete(path)
}



