package firestore_facade

import "cloud.google.com/go/firestore"

// TODO: make an actual facade that simplifies things

// SimpleFirestoreFacade this interface is used instead of full firestore.Client since interfaces should be lean
type SimpleFirestoreFacade interface {
	Doc(string) *firestore.DocumentRef
}
