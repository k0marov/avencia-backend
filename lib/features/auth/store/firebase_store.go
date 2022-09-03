package store

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/features/auth/domain/entities"
)

type FBAuthFacade struct {
	client *auth.Client
}

func NewFBAuthFacade(client *auth.Client) FBAuthFacade {
	return FBAuthFacade{
		client: client,
	}
}

func (a FBAuthFacade) Verify(token string) string {
	info, err := a.client.VerifyIDTokenAndCheckRevoked(context.Background(), token)
	if err != nil {
		return ""
	}
	return info.UID
}

func (a FBAuthFacade) UserByEmail(email string) (entities.User, error) {
	user, err := a.client.GetUserByEmail(context.Background(), email) 
	if err != nil {
		return entities.User{}, core_err.Rethrow("getting user by email from fb", err)
	}
	return entities.User{Id: user.UID}, nil
}
