package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/db"
)

func newDocument(doc *firestore.DocumentSnapshot) db.Document {
	return db.Document{
		Id:        doc.Ref.ID,
		Data:      doc.Data(),
		UpdatedAt: doc.UpdateTime,
		CreatedAt: doc.CreateTime,
	}
}

func newDocuments(docs []*firestore.DocumentSnapshot) (res db.Documents) {
	for _, doc := range docs {
		res = append(res, newDocument(doc))
	}
	return
}
