package db

import (
	"encoding/json"

	"github.com/k0marov/avencia-backend/lib/core/core_err"
)

type JsonDocument struct {
	Path []string
	Data map[string]any
}

type JsonGetter = func(db DB, path []string) (JsonDocument, error)
type JsonSetter = func(db DB, path []string, values map[string]any) error

func JsonGetterImpl(db DB, path []string) (JsonDocument, error) {
	doc, err := db.db.Get(path)
	if err != nil {
		return JsonDocument{}, core_err.Rethrow("getting raw doc", err)
	}
	var jsonDoc map[string]any
	err = json.Unmarshal(doc.Data, &jsonDoc)
	if err != nil {
	  return JsonDocument{}, core_err.Rethrow("unmarshalling json", err)
	}

	return JsonDocument{
		Path: path,
		Data: jsonDoc,
	}, nil
}

func JsonSetterImpl(db DB, path []string, value map[string]any) error {
	valueEnc, err := json.Marshal(value)
	if err != nil {
		return core_err.Rethrow("marshalling value", err)
	}
	return db.db.Set(path, valueEnc)
}

func JsonUpdaterImpl(db DB, path []string, value map[string]any) error {
	return db.db.RunTransaction(func(tDB DB) error {
    doc, err := JsonGetterImpl(tDB, path)
    if err != nil {
    	return core_err.Rethrow("getting current doc", err)
    }
    data := doc.Data
    for k, v := range value {
    	data[k] = v
    }
    return JsonSetterImpl(tDB, path, data)
	})
}
