package db

import (
	"encoding/json"

	"github.com/k0marov/avencia-backend/lib/core/core_err"
)

type JsonDocument struct {
	Path []string
	Data map[string]any
}

type JsonDocuments []JsonDocument

type JsonGetter = func(db DB, path []string) (JsonDocument, error)
type JsonCollectionGetter = func(db DB, path []string) (JsonDocuments, error)
type JsonSetter = func(db DB, path []string, value map[string]any) error
type JsonUpdater = func(db DB, path []string, value map[string]any) error

func unmarshalDoc(d Document) (JsonDocument, error) {
	jsonD := JsonDocument{
		Path: d.Path,
	}
	err := json.Unmarshal(d.Data, &jsonD.Data)
	if err != nil {
		return JsonDocument{}, core_err.Rethrow("unmarshalling json", err)
	}
	return jsonD, nil
}

func JsonGetterImpl(db DB, path []string) (JsonDocument, error) {
	doc, err := db.db.Get(path)
	if err != nil {
		return JsonDocument{}, core_err.Rethrow("getting raw doc", err)
	}
	return unmarshalDoc(doc)

}

func JsonCollectionGetterImpl(db DB, path []string) (JsonDocuments, error) {
	docs, err := db.db.GetCollection(path)
	if err != nil {
		return JsonDocuments{}, core_err.Rethrow("getting raw docs", err)
	}
	jDocs := JsonDocuments{} 
	for _, d := range docs {
		jDoc, err := unmarshalDoc(d)
		if err != nil {
			return JsonDocuments{}, core_err.Rethrow("parsing raw doc", err)
		}
		jDocs = append(jDocs, jDoc)
	}
	return jDocs, nil
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
