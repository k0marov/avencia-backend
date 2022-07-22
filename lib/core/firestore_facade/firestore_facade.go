package firestore_facade

import (
	"cloud.google.com/go/firestore"
	"context"
)

// TODO: make an actual facade that simplifies things

// SimpleFirestoreFacade this interface is used instead of full firestore.Client since interfaces should be lean
type SimpleFirestoreFacade interface {
	Doc(string) *firestore.DocumentRef
}

// TransactionFirestoreFacade this interface is used when a transaction needs to be invoked
type TransactionFirestoreFacade interface {
	SimpleFirestoreFacade
	RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) (err error)
}
