package service

import (
	"context"
	"errors"

	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/store"
)

// UserInfoAdder parses the provided auth header and adds corresponding user info to the context
// if the token is valid. Otherwise, leaves the ctx unchanged.
type UserInfoAdder = func(ctx context.Context, authHeader string) context.Context

func NewUserInfoAdder(verify store.TokenVerifier) UserInfoAdder {
	return func(ctx context.Context, authHeader string) context.Context {
		token := tokenFromHeader(authHeader)
		if token == "" {
			return ctx
		}
		userId := verify(token)
		if userId == "" {
			return ctx
		}

		return context.WithValue(ctx, userContextKey, entities.User{Id: userId})
	}
}

func UserFromCtx(ctx context.Context) (entities.User, error) {
	u, ok := ctx.Value(userContextKey).(entities.User)
	if ok {
		return u, nil
	}

	return entities.User{}, ErrNoUserInContext
}

var ErrNoUserInContext = errors.New("there is no user data in the provided context")

// AddUserToCtx is only for usage in tests
func AddUserToCtx(user entities.User, ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// tokenFromHeader returns "" if header cannot be parsed
func tokenFromHeader(header string) string {
	return header
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)
