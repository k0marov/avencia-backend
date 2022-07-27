package core_err

import (
	"errors"
	"fmt"
)

// ErrNotFound is an error that is passed from store layer to domain layer so that a 404 Client Error is thrown
var ErrNotFound = errors.New("not found")

func Rethrow(description string, err error) error {
	return fmt.Errorf("%s: %w", description, err)
}
