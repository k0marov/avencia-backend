package store

import "github.com/k0marov/avencia-backend/lib/features/auth/domain/entities"

// TokenVerifier should return "" if the provided token is invalid
type TokenVerifier = func(token string) (userId string)

type UserByEmailGetter = func(email string) (entities.User, error)

// AuthFacade Verify returns "" if the provided token is invalid
type AuthFacade interface {
	Verify(token string) (userId string)
	UserByEmail(email string) (entities.User, error)
}
