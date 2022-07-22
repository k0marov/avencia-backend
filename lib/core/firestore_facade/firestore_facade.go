package firestore_facade

import (
	"cloud.google.com/go/firestore"
)

// TODO: make an actual facade that simplifies things

// Simple this interface is used instead of full firestore.Client since interfaces should be lean
type Simple interface {
	Doc(string) *firestore.DocumentRef
	Batch() *firestore.WriteBatch
}

// WriteBatch *firestore.WriteBatch implements this
type WriteBatch interface {
	Set(dr *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) *firestore.WriteBatch
}
