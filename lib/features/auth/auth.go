package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/auth/internal"
)

type User = internal.User

func UserFromCtx(ctx context.Context) (User, error) {
	return internal.UserFromCtx(ctx)
}
func AddUserToCtx(user User, ctx context.Context) context.Context {
	return internal.AddUserToCtx(user, ctx)
}

func NewAuthMiddleware(app *firebase.App) core.Middleware {
	fbAuth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firebase Auth: %v", err)
	}
	return internal.NewFirebaseAuthMiddleware(fbAuth)
}
