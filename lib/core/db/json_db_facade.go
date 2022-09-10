package db

import (
	"encoding/json"

	"github.com/k0marov/avencia-backend/lib/core/core_err"
)

type JsonGetter[T any] func(db DB, path []string) (T, error)
type JsonColGetter[T any] func(db DB, path []string) ([]T, error)
type JsonSetter[T any] func(db DB, path []string, val T) error
type JsonUpdater[T any] func(db DB, path []string, key string, val T) error

func NewJsonGetter[T any]() JsonGetter[T] {
	return func(db DB, path []string) (res T, err error) {
		d, err := db.db.Get(path)
		if err != nil {
			return res, core_err.Rethrow("getting raw doc", err)
		}
		return parseDoc[T](d.Data)
	}
}

func NewJsonColGetter[T any]() JsonColGetter[T] {
	return func(db DB, path []string) (res []T, err error) {
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
}

func NewJsonSetter[T any]() JsonSetter[T] {
	return func(db DB, path []string, val T) error {
		valEncoded, err := json.Marshal(val)
		if err != nil {
			return core_err.Rethrow("marshalling val", err)
		}
		return db.db.Set(path, valEncoded)
	}
}

func NewJsonUpdater[T any](get JsonGetter[map[string]any], set JsonSetter[map[string]any]) JsonUpdater[T] {
	return func(db DB, path []string, key string, val T) error {
    current, err := get(db, path) 
    if err != nil {
    	return core_err.Rethrow("getting current doc", err)
    }
    valJson, err := structToMap(val)
    if err != nil {
    	return core_err.Rethrow("converting val to map", err)
    }
    current[key] = valJson

    return set(db, path, current)
	}
}




func parseDoc[T any](doc []byte) (T, error) {
	var res T
	err := json.Unmarshal(doc, &res)
	if err != nil {
		return res, core_err.Rethrow("unmarshalling doc", err)
	}
	return res, nil
}


// structToMap is currently implemented as a hack of marshalling and then unmarshalling the struct. 
// This can be a performance bottleneck.
func structToMap(s any) (map[string]any, error) {
	var inMap map[string]any
	inJson, err := json.Marshal(s)
	if err != nil {
		return inMap, core_err.Rethrow("marshalling the struct", err)
	}
	err = json.Unmarshal(inJson, &inMap)
	if err != nil {
		return inMap, core_err.Rethrow("unmarshalling the struct", err)
	}
	return inMap, nil
}

