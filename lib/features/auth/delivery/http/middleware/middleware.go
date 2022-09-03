package middleware

import (
	"net/http"

	"github.com/k0marov/avencia-backend/lib/core"
)

func NewAuthMiddleware(addUserInfo service.UserInfoAdder) core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			ctx := addUserInfo(r.Context(), authHeader)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
