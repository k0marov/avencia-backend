package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/k0marov/avencia-backend/lib/core"
)

type Verifier = func(token string) (userId string, ok bool)

func NewFirebaseAuthMiddleware(verify Verifier) core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := tokenFromHeader(r)
			if !ok {
				next.ServeHTTP(w, r)
				return 
			}
			userId, ok := verify(token)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

      newR := r.WithContext(context.WithValue(r.Context(), userContextKey, User{Id: userId}))

			next.ServeHTTP(w, newR)
		})
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

func tokenFromHeader(r *http.Request) (string, bool) {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:], true
	}

	return "", false
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)


