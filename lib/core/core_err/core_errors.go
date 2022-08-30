package core_err

import (
	"fmt"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
)

// ErrNotFound is an error that is passed from store layer to domain layer so that a 404 Client Error is thrown
type errNotFound struct {

}
func (e errNotFound) Error() string {
	return "not found"
}

var ErrNotFound = errNotFound{} 

func IsNotFound(err error) bool {
	_, ok := err.(errNotFound)
	return ok 
}

// TODO: actually this can be a performance bottleneck because of the type conversion
func Rethrow(description string, err error) error {
	clientErr, ok := err.(client_errors.ClientError)
	if ok {
		return clientErr
	}
	return fmt.Errorf("%s: %w", description, err)
}

