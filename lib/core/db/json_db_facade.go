package db

import (
	"encoding/json"

	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
)

type JsonGetter[T any] func(db DB, path []string) (T, error)
type JsonColGetter[T any] func(db DB, path []string) ([]T, error)
type JsonSetter[T any] func(db DB, path []string, val T) error
type JsonUpdater[T any] func(db DB, path []string, key string, val T) error
type JsonMultiUpdater[T any] func(db DB, path []string, newVal T) error

func JsonGetterImpl[T any](db DB, path []string) (res T, err error) {
	d, err := db.db.Get(path)
	if err != nil {
		return res, core_err.Rethrow("getting raw doc", err)
	}
	return parseDoc[T](d.Data)
}

func JsonColGetterImpl[T any](db DB, path []string) (res []T, err error) {
	docs, err := db.db.GetCollection(path)
	if err != nil {
		return res, core_err.Rethrow("getting raw collection elems", err)
	}
	for _, d := range docs {
		parsed, err := parseDoc[T](d.Data)
		if err != nil {
			return res, core_err.Rethrow("parsing one of the raw col docs", err)
		}
		res = append(res, parsed)
	}
	return res, nil
}

func JsonSetterImpl[T any](db DB, path []string, val T) error {
	valEncoded, err := json.Marshal(val)
	if err != nil {
		return core_err.Rethrow("marshalling val", err)
	}
	return db.db.Set(path, valEncoded)
}

func JsonUpdaterImpl[T any](db DB, path []string, key string, val T) error {
	current, err := JsonGetterImpl[map[string]any](db, path)
	if err != nil && !core_err.IsNotFound(err){
		return core_err.Rethrow("getting current doc", err)
	}
	if current == nil {
		current = map[string]any{} 
	}
	current[key] = val

	return JsonSetterImpl(db, path, current)
}

// JsonMultiUpdaterImpl could be a performance bottleneck 
// since newVal is marshalled and then unmarshalled again.
func JsonMultiUpdaterImpl[T any](db DB, path []string, newVal T) error {
	current, err := JsonGetterImpl[map[string]any](db, path)
	if err != nil && !core_err.IsNotFound(err){
		return core_err.Rethrow("getting current doc", err)
	}
	if current == nil {
		current = map[string]any{} 
	}
	newValBin, err := json.Marshal(newVal)
	if err != nil {
		return core_err.Rethrow("marhsalling the new value", err)
	}
	var newValMap map[string]any 
	err = json.Unmarshal(newValBin, &newValMap)
	if err != nil {
		return core_err.Rethrow("while unmarshalling the new value to get a map", err)
	}
	for key, value := range newValMap {
		current[key] = value
	}
	return JsonSetterImpl(db, path, current)
}

func parseDoc[T any](doc *[]byte) (T, error) {
	var res T
	err := json.Unmarshal(*doc, &res)
	if err != nil {
		return res, core_err.Rethrow("unmarshalling doc", err)
	}
	return res, nil
}
