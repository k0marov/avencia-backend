package store

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
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

func (a FBAuthFacade) Get(userId string) (entities.DetailedUser, error) {
	user, err := a.client.GetUser(context.Background(), userId)
	if err != nil {
		return entities.DetailedUser{}, core_err.Rethrow("getting user info from firebase", err)
	}
	return entities.DetailedUser{
		Id:          userId,
		Email:       user.Email,
		PhoneNum:       user.PhoneNumber,
		DisplayName: user.DisplayName,
	}, nil
}
