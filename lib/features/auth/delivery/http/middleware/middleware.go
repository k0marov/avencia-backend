package middleware

import (
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/service"
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
