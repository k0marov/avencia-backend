package auth

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/features/auth/internal"
)

// services

type UserFromEmail = func(email string) (User, error)

func NewUserFromEmail(fbAuth *auth.Client) UserFromEmail {
	return func(email string) (User, error) {
		user, err := fbAuth.GetUserByEmail(context.Background(), email)
		if err != nil {
			return User{}, core_err.Rethrow("getting user from firebase", err)
		}
		return User{Id: user.UID}, nil
	}
}

// user entity

type User = internal.User

func UserFromCtx(ctx context.Context) (User, error) {
	return internal.UserFromCtx(ctx)
}
func AddUserToCtx(user User, ctx context.Context) context.Context {
	return internal.AddUserToCtx(user, ctx)
}

// middleware

func NewAuthMiddleware(fbAuth *auth.Client) core.Middleware {
	return internal.NewFirebaseAuthMiddleware(fbAuth)
}
