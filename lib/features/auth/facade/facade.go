package facade

import (
	"context"

	"firebase.google.com/go/auth"
)

type AuthFacade interface {
	Verify(token string) (userId string, ok bool)
}

type FBAuthFacade struct {
	client *auth.Client
}

func NewFBAuthFacade(client *auth.Client) FBAuthFacade {
  return FBAuthFacade{
    client: client,
  }
}

func (a FBAuthFacade) Verify(token string) (string, bool) {
  info, err := a.client.VerifyIDTokenAndCheckRevoked(context.Background(), token)
  if err != nil {
    return "", false
  }
  return info.UID, true
}
