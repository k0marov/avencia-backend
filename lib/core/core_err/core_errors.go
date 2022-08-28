package core_err

import (
	"errors"
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
)

// ErrNotFound is an error that is passed from store layer to domain layer so that a 404 Client Error is thrown
var ErrNotFound = errors.New("not found")

func Rethrow(description string, err error) error {
	clientErr, ok := err.(client_errors.ClientError)
	if ok {
		return clientErr
	}
	return fmt.Errorf("%s: %w", description, err)
}

