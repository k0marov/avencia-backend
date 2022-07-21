package core_errors

import "errors"

// ErrNotFound is an error that is passed from store layer to domain layer so that a 404 Client Error is thrown
var ErrNotFound = errors.New("not found")
