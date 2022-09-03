package auth

import (
	"context"
	"errors"
	"strings"
)

// UserInfoAdder parses the provided auth header and adds corresponding user info to the context
// if the token is valid. Otherwise, leaves the ctx unchanged.
type UserInfoAdder = func(ctx context.Context, authHeader string) context.Context

// TokenVerifier should return "" if the provided token is invalid
type TokenVerifier = func(token string) (userId string)

func NewUserInfoAdder(verify TokenVerifier) UserInfoAdder {
	return func(ctx context.Context, authHeader string) context.Context {
		token:= tokenFromHeader(authHeader)
		if token == "" {
		  return ctx
		}
		userId := verify(token)
		if userId == "" {
			return ctx
		}

		return context.WithValue(ctx, userContextKey, User{Id: userId})
	}
}

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if ok {
		return u, nil
	}

	return User{}, ErrNoUserInContext
}

var ErrNoUserInContext = errors.New("there is no user data in the provided context")

// AddUserToCtx is only for usage in tests
func AddUserToCtx(user User, ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// tokenFromHeader returns "" if header cannot be parsed
func tokenFromHeader(header string) string {
	if len(header) > 7 && strings.ToLower(header[0:6]) == "bearer" {
		return header[7:]
	}

	return ""
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)
