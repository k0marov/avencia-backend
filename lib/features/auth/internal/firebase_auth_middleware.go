package internal

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/k0marov/avencia-backend/lib/core"
)

func NewFirebaseAuthMiddleware(authClient *auth.Client) core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			bearerToken := tokenFromHeader(r)
			if bearerToken == "" {
				w.WriteHeader(401)
				//httperr.Unauthorised("empty-bearer-token", nil, w, r)
				return
			}

			token, err := authClient.VerifyIDToken(ctx, bearerToken)
			if err != nil {
				w.WriteHeader(401)
				//httperr.Unauthorised("unable-to-verify-jwt", err, w, r)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, userContextKey, User{Id: token.UID})))
		})
	}
}

func tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if ok {
		return u, nil
	}

	return User{}, ErrNoUserInContext
}
