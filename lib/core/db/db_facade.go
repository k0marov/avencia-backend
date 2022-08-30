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
	Get(path []string) (Document, error)
	GetCollection(path []string) (Documents, error)
	Set(path []string, data []byte) error
	RunTransaction(func(DB) error) error
}

type Document struct {
	Path []string
	Data []byte // TODO: maybe store this as a pointer
}

type Documents []Document

type Getter = func(db DB, path []string) (Document, error)
type CollectionGetter = func(db DB, colPath []string) (Documents, error)
type Setter = func(db DB, path []string, data map[string]any) error

func GetterImpl(db DB, path []string) (Document, error) {
	return db.db.Get(path)
}
func CollectionGetterImpl(db DB, colPath []string) (Documents, error) {
	return db.db.GetCollection(colPath)
}

func SetterImpl(db DB, path []string, data []byte) error {
	return db.db.Set(path, data)
}
