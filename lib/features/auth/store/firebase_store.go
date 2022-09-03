package store

import (
	"context"

	"firebase.google.com/go/auth"
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
