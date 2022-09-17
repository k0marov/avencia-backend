package store

import "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"

// TokenVerifier should return "" if the provided token is invalid
type TokenVerifier = func(token string) (userId string)

type UserByEmailGetter = func(email string) (entities.User, error)

type UserGetter = func(userId string) (entities.DetailedUser, error)

// AuthFacade Verify returns "" if the provided token is invalid
type AuthFacade interface {
	Verify(token string) (userId string)
	UserByEmail(email string) (entities.User, error)
	Get(userId string) (entities.DetailedUser, error)
}
