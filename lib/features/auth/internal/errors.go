package internal

import "errors"

var ErrNoUserInContext = errors.New("no user was assigned to this context object")
