package auth

import "errors"

var NoUserInContextErr = errors.New("no user was assigned to this context object")
